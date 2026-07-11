package internal

import "github.com/om-baji/kinetic/shared"

func ValidateDAG(tasks []shared.Task) error {
	if len(tasks) == 0 {
		return shared.NewValidationError("workflow must contain at least one task")
	}

	idSet := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		if t.Id == "" {
			return shared.NewValidationError("task id must not be empty")
		}
		if idSet[t.Id] {
			return shared.NewValidationError("duplicate task id: " + t.Id)
		}
		idSet[t.Id] = true
	}

	for _, t := range tasks {
		for _, dep := range t.Depends {
			if dep == "" {
				return shared.NewValidationError("dependency id must not be empty in task: " + t.Id)
			}
			if !idSet[dep] {
				return shared.NewValidationError("task '" + t.Id + "' depends on non-existent task: " + dep)
			}
		}
	}

	if HasCycle(tasks) {
		return shared.NewValidationError("workflow contains a cycle")
	}

	return nil
}

func HasCycle(tasks []shared.Task) bool {
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
