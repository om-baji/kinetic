package internal

import (
	"context"
	"errors"
	"time"

	"github.com/om-baji/kinetic/shared"
	"gorm.io/gorm"
)

type WorkflowService struct {
	db *gorm.DB
}

func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{db: db}
}

func (s *WorkflowService) Create(payload *shared.Workflow) (*shared.WorkflowRecord, error) {
	if err := ValidateDAG(payload.Tasks); err != nil {
		return nil, err
	}

	ctx := context.Background()

	workflow := &shared.WorkflowRecord{
		Name:   payload.Name,
		Status: shared.WorkflowStatusCreated,
	}

	if err := s.db.WithContext(ctx).Create(workflow).Error; err != nil {
		shared.HandleErr(err)
	}

	taskMap := make(map[string]*shared.TaskRecord, len(payload.Tasks))
	for _, t := range payload.Tasks {
		record := &shared.TaskRecord{
			WorkflowID: workflow.ID,
			Name:       t.Name,
			Status:     shared.TaskStatusPending,
		}
		if err := s.db.WithContext(ctx).Create(record).Error; err != nil {
			shared.HandleErr(err)
		}
		taskMap[t.Id] = record
	}

	for _, t := range payload.Tasks {
		record := taskMap[t.Id]
		for _, depID := range t.Depends {
			dep := &shared.TaskDependency{
				TaskID:          record.ID,
				DependsOnTaskID: taskMap[depID].ID,
			}
			if err := s.db.WithContext(ctx).Create(dep).Error; err != nil {
				shared.HandleErr(err)
			}
		}
	}

	graph := &shared.Graph{}
	if err := s.db.WithContext(ctx).Create(graph).Error; err != nil {
		shared.HandleErr(err)
	}

	workflow.GraphID = graph.ID
	if err := s.db.WithContext(ctx).Save(workflow).Error; err != nil {
		shared.HandleErr(err)
	}

	nodeMap := make(map[string]*shared.GraphNode, len(payload.Tasks))
	for _, t := range payload.Tasks {
		record := taskMap[t.Id]
		node := &shared.GraphNode{
			GraphID: graph.ID,
			TaskID:  record.ID,
			Status:  shared.TaskStatusPending,
		}
		if err := s.db.WithContext(ctx).Create(node).Error; err != nil {
			shared.HandleErr(err)
		}
		nodeMap[t.Id] = node
	}

	for _, t := range payload.Tasks {
		toNode := nodeMap[t.Id]
		for _, depID := range t.Depends {
			fromNode := nodeMap[depID]
			edge := &shared.GraphEdge{
				GraphID:    graph.ID,
				FromNodeID: fromNode.ID,
				ToNodeID:   toNode.ID,
			}
			if err := s.db.WithContext(ctx).Create(edge).Error; err != nil {
				shared.HandleErr(err)
			}
		}
	}

	if err := s.db.WithContext(ctx).Preload("Tasks.Dependencies").Preload("Tasks.Logs").Preload("Graph.Nodes").Preload("Graph.Edges").First(workflow, "id = ?", workflow.ID).Error; err != nil {
		shared.HandleErr(err)
	}

	return workflow, nil
}

func (s *WorkflowService) GetByID(id string) (*shared.WorkflowRecord, error) {
	var record shared.WorkflowRecord
	err := s.db.
		Preload("Tasks.Dependencies").
		Preload("Tasks.Logs").
		Preload("Graph.Nodes").
		Preload("Graph.Edges").
		First(&record, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewNotFoundError("workflow not found: " + id)
		}
		shared.HandleErr(err)
	}
	return &record, nil
}

func (s *WorkflowService) Delete(id string) error {
	record, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if record.Status == shared.WorkflowStatusRunning {
		return shared.NewConflictError("cannot delete a running workflow")
	}

	now := time.Now().UTC()
	for i := range record.Tasks {
		if record.Tasks[i].Status == shared.TaskStatusPending ||
			record.Tasks[i].Status == shared.TaskStatusReady ||
			record.Tasks[i].Status == shared.TaskStatusRetrying {
			s.db.Model(&record.Tasks[i]).Updates(map[string]interface{}{
				"status":     shared.TaskStatusCancelled,
				"updated_at": now,
			})
		}
	}

	s.db.Model(record).Updates(map[string]interface{}{
		"status":     shared.WorkflowStatusCancelled,
		"updated_at": now,
	})

	s.db.Create(&shared.LogEntry{
		Timestamp: now,
		Level:     "info",
		Message:   "workflow cancelled",
	})

	s.db.Delete(record)

	return nil
}

func (s *WorkflowService) Pause(id string) (*shared.WorkflowRecord, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status != shared.WorkflowStatusRunning {
		return nil, shared.NewConflictError("only running workflows can be paused")
	}

	now := time.Now().UTC()
	s.db.Model(record).Updates(map[string]interface{}{
		"status":     shared.WorkflowStatusPaused,
		"updated_at": now,
	})

	s.db.Create(&shared.LogEntry{
		Timestamp: now,
		Level:     "info",
		Message:   "workflow paused",
	})

	record.Status = shared.WorkflowStatusPaused
	record.UpdatedAt = now
	return record, nil
}

func (s *WorkflowService) Resume(id string) (*shared.WorkflowRecord, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status != shared.WorkflowStatusPaused {
		return nil, shared.NewConflictError("only paused workflows can be resumed")
	}

	now := time.Now().UTC()
	s.db.Model(record).Updates(map[string]interface{}{
		"status":     shared.WorkflowStatusRunning,
		"updated_at": now,
	})

	s.db.Create(&shared.LogEntry{
		Timestamp: now,
		Level:     "info",
		Message:   "workflow resumed",
	})

	record.Status = shared.WorkflowStatusRunning
	record.UpdatedAt = now
	return record, nil
}

func (s *WorkflowService) GetGraph(id string) (*shared.Graph, error) {
	record, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	var graph shared.Graph
	err = s.db.
		Preload("Nodes").
		Preload("Edges").
		First(&graph, "id = ?", record.GraphID).Error
	if err != nil {
		shared.HandleErr(err)
	}
	return &graph, nil
}

func (s *WorkflowService) GetLogs(id string) ([]shared.LogEntry, error) {
	_, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	var taskIDs []string
	s.db.Model(&shared.TaskRecord{}).Where("workflow_id = ?", id).Pluck("id", &taskIDs)

	var logs []shared.LogEntry
	if err := s.db.Where("task_id IN ?", taskIDs).Order("timestamp asc").Find(&logs).Error; err != nil {
		shared.HandleErr(err)
	}
	return logs, nil
}
