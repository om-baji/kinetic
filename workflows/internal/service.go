package internal

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/om-baji/kinetic/shared"
)

type WorkflowService struct {
	mu        sync.RWMutex
	workflows map[string]*shared.WorkflowRecord
	logs      map[string][]shared.LogEntry
}

func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		workflows: make(map[string]*shared.WorkflowRecord),
		logs:      make(map[string][]shared.LogEntry),
	}
}

func (s *WorkflowService) Create(payload *shared.Workflow) (*shared.WorkflowRecord, error) {
	if err := ValidateDAG(payload.Tasks); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	id := uuid.New().String()

	taskRecords := make([]shared.TaskRecord, len(payload.Tasks))
	for i, t := range payload.Tasks {
		taskRecords[i] = shared.TaskRecord{
			ID:      t.Id,
			Depends: t.Depends,
			Status:  shared.TaskStatusPending,
		}
	}

	record := &shared.WorkflowRecord{
		ID:        id,
		Name:      payload.Name,
		Status:    shared.WorkflowStatusCreated,
		Tasks:     taskRecords,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mu.Lock()
	s.workflows[id] = record
	s.logs[id] = []shared.LogEntry{
		{
			Timestamp: now,
			TaskID:    "",
			Level:     "info",
			Message:   "workflow created",
		},
	}
	s.mu.Unlock()

	return record, nil
}

func (s *WorkflowService) GetByID(id string) (*shared.WorkflowRecord, error) {
	s.mu.RLock()
	record, ok := s.workflows[id]
	s.mu.RUnlock()

	if !ok {
		return nil, shared.NewNotFoundError("workflow not found: " + id)
	}
	return record, nil
}

func (s *WorkflowService) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.workflows[id]
	if !ok {
		return shared.NewNotFoundError("workflow not found: " + id)
	}

	if record.Status == shared.WorkflowStatusRunning {
		return shared.NewConflictError("cannot delete a running workflow")
	}

	record.Status = shared.WorkflowStatusCancelled
	record.UpdatedAt = time.Now().UTC()

	for i := range record.Tasks {
		if record.Tasks[i].Status == shared.TaskStatusPending ||
			record.Tasks[i].Status == shared.TaskStatusReady ||
			record.Tasks[i].Status == shared.TaskStatusRetrying {
			record.Tasks[i].Status = shared.TaskStatusCancelled
		}
	}

	s.appendLogLocked(id, "", "info", "workflow cancelled")

	delete(s.workflows, id)
	return nil
}

func (s *WorkflowService) Pause(id string) (*shared.WorkflowRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.workflows[id]
	if !ok {
		return nil, shared.NewNotFoundError("workflow not found: " + id)
	}

	if record.Status != shared.WorkflowStatusRunning {
		return nil, shared.NewConflictError("only running workflows can be paused")
	}

	record.Status = shared.WorkflowStatusPaused
	record.UpdatedAt = time.Now().UTC()
	s.appendLogLocked(id, "", "info", "workflow paused")

	return record, nil
}

func (s *WorkflowService) Resume(id string) (*shared.WorkflowRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.workflows[id]
	if !ok {
		return nil, shared.NewNotFoundError("workflow not found: " + id)
	}

	if record.Status != shared.WorkflowStatusPaused {
		return nil, shared.NewConflictError("only paused workflows can be resumed")
	}

	record.Status = shared.WorkflowStatusRunning
	record.UpdatedAt = time.Now().UTC()
	s.appendLogLocked(id, "", "info", "workflow resumed")

	return record, nil
}

func (s *WorkflowService) GetGraph(id string) (*shared.Graph, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.workflows[id]
	if !ok {
		return nil, shared.NewNotFoundError("workflow not found: " + id)
	}

	graph := &shared.Graph{
		Nodes: make([]shared.GraphNode, 0, len(record.Tasks)),
		Edges: make([]shared.GraphEdge, 0),
	}

	for _, t := range record.Tasks {
		graph.Nodes = append(graph.Nodes, shared.GraphNode{
			ID:     t.ID,
			Status: t.Status,
		})
		for _, dep := range t.Depends {
			graph.Edges = append(graph.Edges, shared.GraphEdge{
				From: dep,
				To:   t.ID,
			})
		}
	}

	return graph, nil
}

func (s *WorkflowService) GetLogs(id string) ([]shared.LogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.workflows[id]
	if !ok {
		return nil, shared.NewNotFoundError("workflow not found: " + id)
	}

	logs, exists := s.logs[id]
	if !exists {
		return []shared.LogEntry{}, nil
	}

	result := make([]shared.LogEntry, len(logs))
	copy(result, logs)
	return result, nil
}

func (s *WorkflowService) appendLogLocked(workflowID, taskID, level, message string) {
	s.logs[workflowID] = append(s.logs[workflowID], shared.LogEntry{
		Timestamp: time.Now().UTC(),
		TaskID:    taskID,
		Level:     level,
		Message:   message,
	})
}
