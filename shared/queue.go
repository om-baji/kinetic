package shared

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

func EnqueWorkflow(task string, body WorkflowRecord) (*asynq.Task, error) {
	payload, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(task, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

func DequeWorkflow(ctx context.Context, t *asynq.Task) (*WorkflowRecord, error) {
	var wr *WorkflowRecord

	err := json.Unmarshal(t.Payload(), &wr)

	if err != nil {
		return nil, err
	}

	return wr, nil
}
