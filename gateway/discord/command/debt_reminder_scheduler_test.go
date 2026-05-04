package command

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeJob struct {
	id uuid.UUID
}

func (j *fakeJob) ID() uuid.UUID                     { return j.id }
func (j *fakeJob) LastRun() (time.Time, error)       { return time.Time{}, nil }
func (j *fakeJob) Name() string                      { return "fake" }
func (j *fakeJob) NextRun() (time.Time, error)       { return time.Time{}, nil }
func (j *fakeJob) NextRuns(int) ([]time.Time, error) { return nil, nil }
func (j *fakeJob) RunNow() error                     { return nil }
func (j *fakeJob) Tags() []string                    { return nil }

type fakeScheduler struct {
	mu             sync.Mutex
	newJobCalls    int
	newJobErr      error
	removeJobCalls []uuid.UUID
	removeJobErr   error
	nextID         uuid.UUID
}

func (f *fakeScheduler) NewJob(
	_ gocron.JobDefinition, _ gocron.Task, _ ...gocron.JobOption,
) (gocron.Job, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.newJobCalls++

	if f.newJobErr != nil {
		return nil, f.newJobErr
	}

	id := f.nextID
	if id == uuid.Nil {
		id = uuid.New()
	}

	f.nextID = uuid.Nil

	return &fakeJob{id: id}, nil
}

func (f *fakeScheduler) RemoveJob(id uuid.UUID) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.removeJobCalls = append(f.removeJobCalls, id)

	return f.removeJobErr
}

func newTestScheduler(fake *fakeScheduler) *DebtReminderScheduler {
	return &DebtReminderScheduler{scheduler: fake}
}

func TestScheduleProductionRun_FirstCallNoReplace(t *testing.T) {
	fake := &fakeScheduler{nextID: uuid.New()}
	r := newTestScheduler(fake)

	runAt := time.Now().Add(15 * 24 * time.Hour)

	result, err := r.ScheduleProductionRun(runAt, func() {})
	require.NoError(t, err)
	require.NoError(t, result.RemoveWarn)
	require.True(t, result.ReplacedAt.IsZero(), "no prior schedule so ReplacedAt must be zero")
	require.Equal(t, 1, fake.newJobCalls)
	require.Empty(t, fake.removeJobCalls)
	require.NotEqual(t, uuid.Nil, r.pending)
	require.Equal(t, runAt, r.pendingAt)
}

func TestScheduleProductionRun_SecondCallReplacesFirst(t *testing.T) {
	firstID := uuid.New()
	secondID := uuid.New()

	fake := &fakeScheduler{nextID: firstID}
	r := newTestScheduler(fake)

	firstRunAt := time.Now().Add(15 * 24 * time.Hour)

	_, err := r.ScheduleProductionRun(firstRunAt, func() {})
	require.NoError(t, err)

	fake.nextID = secondID
	secondRunAt := firstRunAt.Add(24 * time.Hour)

	result, err := r.ScheduleProductionRun(secondRunAt, func() {})
	require.NoError(t, err)
	require.NoError(t, result.RemoveWarn)
	require.Equal(t, firstRunAt, result.ReplacedAt)
	require.Equal(t, 2, fake.newJobCalls)
	require.Equal(t, []uuid.UUID{firstID}, fake.removeJobCalls)
	require.Equal(t, secondID, r.pending)
	require.Equal(t, secondRunAt, r.pendingAt)
}

func TestScheduleProductionRun_RemoveErrorDoesNotBlockSchedule(t *testing.T) {
	firstID := uuid.New()
	secondID := uuid.New()

	fake := &fakeScheduler{nextID: firstID}
	r := newTestScheduler(fake)

	firstRunAt := time.Now().Add(15 * 24 * time.Hour)

	_, err := r.ScheduleProductionRun(firstRunAt, func() {})
	require.NoError(t, err)

	removeErr := errors.New("boom")
	fake.removeJobErr = removeErr
	fake.nextID = secondID
	secondRunAt := firstRunAt.Add(24 * time.Hour)

	result, err := r.ScheduleProductionRun(secondRunAt, func() {})
	require.NoError(t, err)
	require.ErrorIs(t, result.RemoveWarn, removeErr)
	require.Equal(t, firstRunAt, result.ReplacedAt)
	require.Equal(t, secondID, r.pending)
}

func TestScheduleProductionRun_NewJobErrorLeavesNoPending(t *testing.T) {
	firstID := uuid.New()

	fake := &fakeScheduler{nextID: firstID}
	r := newTestScheduler(fake)

	firstRunAt := time.Now().Add(15 * 24 * time.Hour)

	_, err := r.ScheduleProductionRun(firstRunAt, func() {})
	require.NoError(t, err)

	newErr := errors.New("new job failed")
	fake.newJobErr = newErr

	secondRunAt := firstRunAt.Add(24 * time.Hour)

	_, err = r.ScheduleProductionRun(secondRunAt, func() {})
	require.Error(t, err)
	require.ErrorIs(t, err, newErr)

	require.Equal(t, uuid.Nil, r.pending)
	require.True(t, r.pendingAt.IsZero())
}

func TestScheduleProductionRun_ConcurrentCallsSerializeCorrectly(t *testing.T) {
	fake := &fakeScheduler{}
	r := newTestScheduler(fake)

	const goroutines = 10

	var wg sync.WaitGroup

	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			_, err := r.ScheduleProductionRun(time.Now().Add(time.Hour), func() {})
			require.NoError(t, err)
		}()
	}

	wg.Wait()

	require.Equal(t, goroutines, fake.newJobCalls)
	require.Len(t, fake.removeJobCalls, goroutines-1, "every call after the first should trigger a remove")
	require.NotEqual(t, uuid.Nil, r.pending)
}

func TestClearIfMatches_OnlyClearsOnMatch(t *testing.T) {
	fake := &fakeScheduler{}
	r := newTestScheduler(fake)

	pending := uuid.New()
	r.pending = pending
	r.pendingAt = time.Now()

	r.clearIfMatches(uuid.New())
	require.Equal(t, pending, r.pending, "non-matching id must not clear")
	require.False(t, r.pendingAt.IsZero())

	r.clearIfMatches(pending)
	require.Equal(t, uuid.Nil, r.pending)
	require.True(t, r.pendingAt.IsZero())
}
