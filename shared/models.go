package shared

import "time"

type Workflow struct {
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id      string   `json:"id"`
	Depends []string `json:"depends"`
}

type WorkflowStatus string

const (
	WorkflowStatusCreated   WorkflowStatus = "created"
	WorkflowStatusQueued    WorkflowStatus = "queued"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusPaused    WorkflowStatus = "paused"
	WorkflowStatusCompleted WorkflowStatus = "completed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
	WorkflowStatusFailed    WorkflowStatus = "failed"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusReady     TaskStatus = "ready"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusSucceeded TaskStatus = "succeeded"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
	TaskStatusRetrying  TaskStatus = "retrying"
	TaskStatusTimedOut  TaskStatus = "timed_out"
)

type WorkflowRecord struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Status    WorkflowStatus `json:"status"`
	Tasks     []TaskRecord   `json:"tasks"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TaskRecord struct {
	ID      string     `json:"id"`
	Depends []string   `json:"depends"`
	Status  TaskStatus `json:"status"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GraphNode struct {
	ID     string     `json:"id"`
	Status TaskStatus `json:"status"`
}

type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	TaskID    string    `json:"task_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
