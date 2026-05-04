package command

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBuildProductionResponse_NoReplaceNoWarn(t *testing.T) {
	runAt := time.Date(2026, 5, 5, 9, 0, 0, 0, time.UTC)

	got := buildProductionResponse(runAt, ScheduleResult{})
	require.Equal(
		t,
		"提醒已執行（模式: 正式）。下次執行: 2026-05-05 09:00",
		got,
	)
}

func TestBuildProductionResponse_ReplacedPrior(t *testing.T) {
	runAt := time.Date(2026, 5, 5, 9, 0, 0, 0, time.UTC)
	replacedAt := time.Date(2026, 4, 20, 14, 45, 0, 0, time.UTC)

	got := buildProductionResponse(runAt, ScheduleResult{ReplacedAt: replacedAt})
	require.Equal(
		t,
		"提醒已執行（模式: 正式）。下次執行: 2026-05-05 09:00（已取代先前排程: 2026-04-20 14:45）",
		got,
	)
}

func TestBuildProductionResponse_ReplacedWithRemoveWarning(t *testing.T) {
	runAt := time.Date(2026, 5, 5, 9, 0, 0, 0, time.UTC)
	replacedAt := time.Date(2026, 4, 20, 14, 45, 0, 0, time.UTC)
	warn := errors.New("boom")

	got := buildProductionResponse(runAt, ScheduleResult{ReplacedAt: replacedAt, RemoveWarn: warn})
	require.Equal(
		t,
		"提醒已執行（模式: 正式）。下次執行: 2026-05-05 09:00（已取代先前排程: 2026-04-20 14:45）（注意：先前排程移除失敗: boom）",
		got,
	)
}
