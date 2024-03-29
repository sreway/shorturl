package shortener

import (
	"context"
	"time"

	"golang.org/x/exp/slog"

	"github.com/sreway/shorturl/internal/domain/url"
)

type (
	action string
	// task implements the task type for the queue.
	task struct {
		name action
		urls []url.URL
	}
)

const (
	deleteAction action = "delete"
)

// NewTask implements the creation of task for queue.
func NewTask(name action, urls []url.URL) *task {
	return &task{
		name: name,
		urls: urls,
	}
}

// ProcQueue implements processing task queue.
func (uc *useCase) ProcQueue(ctx context.Context, checkInterval time.Duration) error {
	tick := time.NewTicker(checkInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if len(uc.taskQueue) == 0 {
				continue
			}

			actions := make(map[action][]url.URL, len(uc.taskQueue))

			for len(uc.taskQueue) != 0 {
				t := <-uc.taskQueue
				actions[t.name] = append(actions[t.name], t.urls...)
			}

			for k, v := range actions {
				switch k {
				case deleteAction:
					err := uc.storage.BatchDelete(ctx, v)
					if err != nil {
						uc.logger.Error("failed batch update", err, slog.String("func", "ProcQueue"))
						continue
					}
				default:
					uc.logger.Warn("unknown task action", slog.Any("action", k),
						slog.String("func", "ProcQueue"))
				}
			}
		case <-ctx.Done():
			close(uc.taskQueue)
			uc.logger.Info("stop processed task queue")
			return nil
		}
	}
}
