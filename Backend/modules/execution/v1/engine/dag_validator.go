package engine

import "fmt"

func Validate(d *DAGDefinition) error {
	if d == nil {
		return fmt.Errorf("dag definition is required")
	}
	if len(d.Steps) == 0 {
		return fmt.Errorf("at least one step is required")
	}
	steps := map[string]Step{}
	for _, s := range d.Steps {
		if s.ID == "" || s.Name == "" || s.Type == "" {
			return fmt.Errorf("step id, name, and type are required")
		}
		if _, ok := steps[s.ID]; ok {
			return fmt.Errorf("duplicate step id: %s", s.ID)
		}
		switch s.Type {
		case StepTypeHTTPCall, StepTypeScript, StepTypeDelay, StepTypeBranch:
		default:
			return fmt.Errorf("unsupported step type: %s", s.Type)
		}
		steps[s.ID] = s
	}
	adj := map[string][]string{}
	for _, s := range d.Steps {
		for _, dep := range s.DependsOn {
			if _, ok := steps[dep]; !ok {
				return fmt.Errorf("unknown dependency %s for step %s", dep, s.ID)
			}
			adj[dep] = append(adj[dep], s.ID)
		}
	}
	for _, e := range d.Edges {
		if _, ok := steps[e.From]; !ok {
			return fmt.Errorf("unknown edge from: %s", e.From)
		}
		if _, ok := steps[e.To]; !ok {
			return fmt.Errorf("unknown edge to: %s", e.To)
		}
		adj[e.From] = append(adj[e.From], e.To)
	}
	visiting := map[string]bool{}
	visited := map[string]bool{}
	var dfs func(string) error
	dfs = func(id string) error {
		if visiting[id] {
			return fmt.Errorf("cycle detected at step %s", id)
		}
		if visited[id] {
			return nil
		}
		visiting[id] = true
		for _, next := range adj[id] {
			if err := dfs(next); err != nil {
				return err
			}
		}
		visiting[id] = false
		visited[id] = true
		return nil
	}
	for id := range steps {
		if err := dfs(id); err != nil {
			return err
		}
	}
	return nil
}
