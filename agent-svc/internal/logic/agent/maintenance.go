package agent

import (
	"context"
	"strings"
	"time"

	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	runtimeMaintenanceInterval = 30 * time.Second
	processingRunTimeout       = 2 * time.Minute
	generatingRunTimeout       = 45 * time.Second
)

func (s *sAgent) startRuntimeMaintenance() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Error(context.Background(), "agent runtime maintenance panic recovered:", r)
			}
		}()

		ticker := time.NewTicker(runtimeMaintenanceInterval)
		defer ticker.Stop()

		for range ticker.C {
			func() {
				defer func() {
					if r := recover(); r != nil {
						g.Log().Error(context.Background(), "agent runtime maintenance tick panic recovered:", r)
					}
				}()

				if err := s.reapZombieRuns(context.Background()); err != nil {
					g.Log().Warning(context.Background(), "agent runtime maintenance failed:", err)
				}
			}()
		}
	}()
}

func (s *sAgent) reapZombieRuns(ctx context.Context) error {
	if err := s.failExpiredToolRuns(ctx); err != nil {
		return err
	}
	if err := s.failStaleRuns(ctx, runtimeStatusGenerating, gtime.Now().Add(-generatingRunTimeout), "GENERATING_TIMEOUT", "run stayed in generating for too long"); err != nil {
		return err
	}
	if err := s.failStaleRuns(ctx, runtimeStatusProcessing, gtime.Now().Add(-processingRunTimeout), "PROCESSING_STALE", "run was orphaned during processing"); err != nil {
		return err
	}
	return nil
}

func (s *sAgent) failExpiredToolRuns(ctx context.Context) error {
	var rows []entity.AgentRun
	err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().RunStatusCode, runtimeStatusWaitingToolCall).
		WhereLTE(dao.AgentRun.Columns().ToolWaitTimeoutAt, gtime.Now()).
		Scan(&rows)
	if err != nil {
		return gerror.Wrap(err, "query expired tool wait runs failed")
	}
	for i := range rows {
		if err = s.failRun(ctx, &rows[i], "TOOL_WAIT_TIMEOUT", "tool call wait timed out"); err != nil {
			return err
		}
	}
	return nil
}

func (s *sAgent) failStaleRuns(ctx context.Context, runStatus string, cutoff *gtime.Time, errorCode, errorMessage string) error {
	var rows []entity.AgentRun
	err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().RunStatusCode, runStatus).
		WhereLTE(dao.AgentRun.Columns().UpdatedAt, cutoff).
		Scan(&rows)
	if err != nil {
		return gerror.Wrap(err, "query stale runs failed")
	}
	for i := range rows {
		if err = s.failRun(ctx, &rows[i], errorCode, errorMessage); err != nil {
			return err
		}
	}
	return nil
}

func (s *sAgent) failRun(ctx context.Context, run *entity.AgentRun, errorCode, errorMessage string) error {
	if run == nil || run.Id == 0 {
		return nil
	}

	now := gtime.Now()
	_, err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().RunNo, run.RunNo).
		Where(dao.AgentRun.Columns().RunStatusCode, run.RunStatusCode).
		Data(do.AgentRun{
			RunStatusCode:     runtimeStatusFailed,
			CurrentNodeCode:   firstNonEmpty(run.CurrentNodeCode, "RuntimeMaintenance"),
			ToolResultStatus:  normalizeToolResultStatus(run.ToolResultStatus),
			ErrorCode:         strings.TrimSpace(errorCode),
			ErrorMessage:      strings.TrimSpace(errorMessage),
			CheckpointVersion: gdb.Raw("checkpoint_version + 1"),
			FinishedAt:        now,
			UpdatedAt:         now,
		}).
		Update()
	if err != nil {
		return gerror.Wrap(err, "fail stale run failed")
	}

	_, _ = dao.AgentConversation.Ctx(ctx).
		Where(dao.AgentConversation.Columns().ConversationNo, run.ConversationNo).
		Data(do.AgentConversation{
			LastRunStatusCode: runtimeStatusFailed,
			UpdatedAt:         now,
		}).
		Update()

	return nil
}
