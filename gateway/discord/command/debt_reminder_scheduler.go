package command

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// schedulerAPI is the subset of gocron.Scheduler used by DebtReminderScheduler.
// Narrowing the interface keeps unit tests hermetic.
type schedulerAPI interface {
	NewJob(def gocron.JobDefinition, task gocron.Task, opts ...gocron.JobOption) (gocron.Job, error)
	RemoveJob(id uuid.UUID) error
}

// DebtReminderScheduler wraps a gocron scheduler and enforces the
// "at most one pending production OneTimeJob" invariant.
type DebtReminderScheduler struct {
	scheduler schedulerAPI

	mu        sync.Mutex
	pending   uuid.UUID
	pendingAt time.Time
}

// ScheduleResult reports how a ScheduleProductionRun call interacted with any
// previously pending job. ReplacedAt is zero if nothing was replaced.
// RemoveWarn is non-nil only when removing the prior pending job failed;
// it is a warning rather than a fatal error — the new job was still added.
type ScheduleResult struct {
	ReplacedAt time.Time
	RemoveWarn error
}

func NewDebtReminderScheduler(s gocron.Scheduler) *DebtReminderScheduler {
	return &DebtReminderScheduler{scheduler: s}
}

// ScheduleProductionRun replaces any currently pending production OneTimeJob
// with a new one that fires at runAt. The provided task is invoked on
// execution; after it returns, the stored UUID is cleared iff it still
// matches (avoiding clobbering a newer replacement that raced in).
//
// If NewJob fails, the prior pending state has already been cleared (the
// remove call succeeded) so no pending job exists afterwards. The returned
// ScheduleResult still reports the ReplacedAt of the cleared prior job.
func (r *DebtReminderScheduler) ScheduleProductionRun(
	runAt time.Time, task func(),
) (ScheduleResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := ScheduleResult{}

	if r.pending != uuid.Nil {
		rmErr := r.scheduler.RemoveJob(r.pending)
		if rmErr != nil {
			log.Printf("debt-reminder: remove prior pending job failed: %s", rmErr)
			result.RemoveWarn = rmErr
		}

		result.ReplacedAt = r.pendingAt
		r.pending = uuid.Nil
		r.pendingAt = time.Time{}
	}

	var jobID uuid.UUID

	job, err := r.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(runAt)),
		gocron.NewTask(func() {
			task()
			r.clearIfMatches(jobID)
		}),
	)
	if err != nil {
		return result, fmt.Errorf("schedule production run: %w", err)
	}

	jobID = job.ID()
	r.pending = jobID
	r.pendingAt = runAt

	return result, nil
}

func (r *DebtReminderScheduler) clearIfMatches(id uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pending == id {
		r.pending = uuid.Nil
		r.pendingAt = time.Time{}
	}
}
