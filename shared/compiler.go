package shared

import (
	"context"

	"gorm.io/gorm"
)

func ValidateDAG(tasks []Task) error {
	if len(tasks) == 0 {
		return NewValidationError("workflow must contain at least one task")
	}

	idSet := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		if t.Id == "" {
			return NewValidationError("task id must not be empty")
		}
		if idSet[t.Id] {
			return NewValidationError("duplicate task id: " + t.Id)
		}
		idSet[t.Id] = true
	}

	for _, t := range tasks {
		for _, dep := range t.Depends {
			if dep == "" {
				return NewValidationError("dependency id must not be empty in task: " + t.Id)
			}
			if !idSet[dep] {
				return NewValidationError("task '" + t.Id + "' depends on non-existent task: " + dep)
			}
		}
	}

	if HasCycle(tasks) {
		return NewValidationError("workflow contains a cycle")
	}

	return nil
}

func HasCycle(tasks []Task) bool {
	adj := make(map[string][]string, len(tasks))
	for _, t := range tasks {
		adj[t.Id] = t.Depends
	}

	const (
		white = 0
		gray  = 1
		black = 2
	)

	color := make(map[string]int, len(tasks))

	var dfs func(node string) bool
	dfs = func(node string) bool {
		color[node] = gray
		for _, neighbor := range adj[node] {
			if color[neighbor] == gray {
				return true
			}
			if color[neighbor] == white {
				if dfs(neighbor) {
					return true
				}
			}
		}
		color[node] = black
		return false
	}

	for _, t := range tasks {
		if color[t.Id] == white {
			if dfs(t.Id) {
				return true
			}
		}
	}

	return false
}

func CompileWorkflow(db gorm.DB, workflow *WorkflowRecord, tasks []Task) (*WorkflowRecord, error) {
	ctx := context.Background()

	taskMap := make(map[string]*TaskRecord, len(tasks))
	for _, t := range tasks {
		record := &TaskRecord{
			WorkflowID: workflow.ID,
			Name:       t.Name,
			Status:     TaskStatusPending,
		}
		if err := db.WithContext(ctx).Create(record).Error; err != nil {
			HandleErr(err)
		}
		taskMap[t.Id] = record
	}

	for _, t := range tasks {
		record := taskMap[t.Id]
		for _, depID := range t.Depends {
			dep := &TaskDependency{
				TaskID:          record.ID,
				DependsOnTaskID: taskMap[depID].ID,
			}
			if err := db.WithContext(ctx).Create(dep).Error; err != nil {
				HandleErr(err)
			}
		}
	}

	graph := &Graph{}
	if err := db.WithContext(ctx).Create(graph).Error; err != nil {
		HandleErr(err)
	}

	workflow.GraphID = graph.ID
	if err := db.WithContext(ctx).Save(workflow).Error; err != nil {
		HandleErr(err)
	}

	nodeMap := make(map[string]*GraphNode, len(tasks))
	for _, t := range tasks {
		record := taskMap[t.Id]
		node := &GraphNode{
			GraphID: graph.ID,
			TaskID:  record.ID,
			Status:  TaskStatusPending,
		}
		if err := db.WithContext(ctx).Create(node).Error; err != nil {
			HandleErr(err)
		}
		nodeMap[t.Id] = node
	}

	for _, t := range tasks {
		toNode := nodeMap[t.Id]
		for _, depID := range t.Depends {
			fromNode := nodeMap[depID]
			edge := &GraphEdge{
				GraphID:    graph.ID,
				FromNodeID: fromNode.ID,
				ToNodeID:   toNode.ID,
			}
			if err := db.WithContext(ctx).Create(edge).Error; err != nil {
				HandleErr(err)
			}
		}
	}

	if err := db.WithContext(ctx).Preload("Tasks.Dependencies").Preload("Tasks.Logs").Preload("Graph.Nodes").Preload("Graph.Edges").First(workflow, "id = ?", workflow.ID).Error; err != nil {
		HandleErr(err)
	}

	return workflow, nil
}
