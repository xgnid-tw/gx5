# /debt-reminder Scheduling Rules

Date: 2026-04-20
Status: Approved, ready for implementation
Target branch: `feat/debt-reminder-command`

## Background

The current `/debt-reminder` handler (`gateway/discord/command/debt_reminder_handler.go`) has two defects observed in production:

1. **Debug invocations still schedule a real production OneTimeJob.** The immediate run respects the `debug` flag, but the delayed task closure always calls `uc.Execute(ctx, false)`. A debug test today leaves behind a real production job N days later.
2. **No deduplication.** Each invocation appends a new OneTimeJob. Repeated invocations stack, so the use case fires multiple times.

Observed symptom: three "debt-reminder scheduled run" log entries on 2026-04-20 at 14:45:57, 14:48:16, 14:51:40 — matching three `/debt-reminder` invocations on 2026-04-05.

## Goals

- Debug invocations are side-effect free on the scheduler.
- At most one pending production OneTimeJob exists at any time.
- The scheduled task always runs in production mode.

## Non-goals

- Persistence across bot restarts. OneTimeJobs remain in-memory-only.
- Cancelling a pending job via a dedicated command. Out of scope here.

## Behavior

```
/debt-reminder days=N debug=B invoked
        │
        ├── B == true:
        │     uc.Execute(ctx, true)   // log-only notification
        │     respond: "提醒已執行（模式: 除錯）。未排程下次執行。"
        │     DONE (scheduler untouched)
        │
        └── B == false:
              uc.Execute(ctx, false)  // real DMs + log
              scheduler.ScheduleProductionRun(now+N*24h, task)
                  ├── if pending UUID stored:
                  │       scheduler.RemoveJob(pending)   // best-effort
                  ├── add new OneTimeJob at runAt
                  └── store its UUID
              respond: "提醒已執行（模式: 正式）。下次執行: YYYY-MM-DD HH:mm
                       [ （已取代先前排程: YYYY-MM-DD HH:mm） if replaced ]"
```

Invariants:

- At most one pending production OneTimeJob exists at any time.
- Debug invocations never create, remove, or inspect a pending job.
- The scheduled OneTimeJob's task always calls `uc.Execute(ctx, false)`, regardless of the invocation that scheduled it.

## Components

### New file: `gateway/discord/command/debt_reminder_scheduler.go`

```go
type DebtReminderScheduler struct {
    scheduler schedulerAPI
    mu        sync.Mutex
    pending   uuid.UUID   // uuid.Nil = none
    pendingAt time.Time
}

func NewDebtReminderScheduler(s gocron.Scheduler) *DebtReminderScheduler

// Returns replacedAt (zero if nothing was replaced) and any error from NewJob.
// If removing the prior pending job fails, logs a warning and continues.
func (r *DebtReminderScheduler) ScheduleProductionRun(
    runAt time.Time, task func(),
) (replacedAt time.Time, err error)

// Called by the task closure on execution to clear the stored UUID
// only if it still matches (avoids clobbering a newer replacement).
func (r *DebtReminderScheduler) clearIfMatches(id uuid.UUID)
```

Internal `schedulerAPI` interface holds only the two gocron methods used, enabling hermetic unit tests:

```go
type schedulerAPI interface {
    NewJob(gocron.JobDefinition, gocron.Task, ...gocron.JobOption) (gocron.Job, error)
    RemoveJob(uuid.UUID) error
}
```

### Changed: `gateway/discord/command/debt_reminder_handler.go`

- `RegisterDebtReminderCommand` takes `*DebtReminderScheduler` instead of `gocron.Scheduler`.
- `handleDebtReminder` splits on `debug`:
  - `debug == true` → `uc.Execute(ctx, true)` → respond "未排程下次執行" → return.
  - `debug == false` → `uc.Execute(ctx, false)` → `sched.ScheduleProductionRun(runAt, taskFn)` → respond.

### Changed: `main.go`

```go
sched := discordcmd.NewDebtReminderScheduler(s)
discordcmd.RegisterDebtReminderCommand(cmdHandler, notifyUnpaidUC, sched)
```

## Error handling

| Failure point | Behavior | User-facing response |
|---|---|---|
| `uc.Execute` fails (immediate, debug) | log error; skip scheduling | `"提醒執行失敗: <err>"` |
| `uc.Execute` fails (immediate, prod) | log error; **do NOT schedule** next run | `"提醒執行失敗，未排程下次執行: <err>"` |
| `RemoveJob(pending)` fails | log warning; **continue** adding new job; clear stored UUID | response includes `"（注意：先前排程移除失敗: <err>）"` |
| `NewJob` fails | return error; stored UUID unchanged (prior pending remains valid) | `"提醒已執行，但排程失敗: <err>"` |

Rationale for continuing after `RemoveJob` fails: gocron's `RemoveJob` can fail if the job already fired or was removed concurrently. The desired invariant is "one pending going forward", not "perfect cleanup of the past". Log, continue.

Scheduled OneTimeJob closure clears the stored UUID via `clearIfMatches` to avoid clobbering a newer replacement that raced in.

## Tests

### `debt_reminder_scheduler_test.go`

Uses a fake `schedulerAPI` implementation (not real gocron).

1. First production schedule: `NewJob` called once, `RemoveJob` never called, UUID stored, `replacedAt` is zero.
2. Second production schedule replaces first: `RemoveJob` called with first UUID; `NewJob` called; stored UUID updated; `replacedAt == firstRunAt`.
3. `RemoveJob` returns error: `NewJob` still called; stored UUID updated; error surfaced via returned warning.
4. `NewJob` returns error: stored UUID unchanged; error returned.
5. Concurrent schedule calls (10 goroutines, `-race`): exactly one UUID remains pending; no data races.
6. `clearIfMatches`: only clears on exact match.

### `debt_reminder_handler_test.go`

Uses fake `port.DebtReminder` and a scheduler interface the handler depends on.

1. `debug=true`: `uc.Execute` called with `true`; scheduler NOT called; response contains "模式: 除錯" and "未排程下次執行".
2. `debug=false`: `uc.Execute` called with `false`; scheduler called with `runAt ≈ now + days*24h`; response contains next-run timestamp.
3. `debug=false` with prior pending: response contains "已取代先前排程".
4. `debug=false`, immediate `uc.Execute` fails: scheduler NOT called; response is the failure message.
5. `days < 1`: existing validation triggers; no immediate run; no scheduling.

## Rollout

- PR targets `feat/debt-reminder-command` (not `main`).
- After merge, deploy via existing CircleCI pipeline (which deploys `main` — see Open Questions).

## Open questions

- None for the scheduling rules themselves. The CI deploy step is tied to `main`; the feat branch's deploy path is out of scope for this change.
