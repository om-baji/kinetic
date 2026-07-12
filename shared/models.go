package shared

import (
	"time"
)

type Workflow struct {
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
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
	BaseModel
	Name   string         `json:"name"`
	Status WorkflowStatus `json:"status"`

	GraphID string
	Graph   Graph

	Tasks []TaskRecord `gorm:"foreignKey:WorkflowID"`
}

type TaskRecord struct {
	BaseModel
	WorkflowID string

	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`

	Dependencies []TaskDependency `gorm:"foreignKey:TaskID"`
	Logs         []LogEntry       `gorm:"foreignKey:TaskID"`
}

type TaskDependency struct {
	BaseModel
	TaskID          string
	DependsOnTaskID string

	Task      TaskRecord `gorm:"foreignKey:TaskID"`
	DependsOn TaskRecord `gorm:"foreignKey:DependsOnTaskID"`
}

type Graph struct {
	BaseModel
	Nodes []GraphNode `gorm:"foreignKey:GraphID"`
	Edges []GraphEdge `gorm:"foreignKey:GraphID"`
}

type GraphNode struct {
	BaseModel
	GraphID string

	TaskID string
	Task   TaskRecord

	Status TaskStatus `json:"status"`
}

type GraphEdge struct {
	BaseModel
	GraphID string

	FromNodeID string
	ToNodeID   string

	FromNode GraphNode `gorm:"foreignKey:FromNodeID"`
	ToNode   GraphNode `gorm:"foreignKey:ToNodeID"`
}

type LogEntry struct {
	BaseModel
	TaskID string

	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
