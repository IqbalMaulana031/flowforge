package engine

import "fmt"

func TopologicalWaves(d *DAGDefinition) ([][]Step, error) {
	if err := Validate(d); err != nil {
		return nil, err
	}
	steps := map[string]Step{}
	indeg := map[string]int{}
	adj := map[string][]string{}
	for _, s := range d.Steps {
		steps[s.ID] = s
		indeg[s.ID] = 0
	}
	for _, s := range d.Steps {
		for _, dep := range s.DependsOn {
			adj[dep] = append(adj[dep], s.ID)
			indeg[s.ID]++
		}
	}
	for _, e := range d.Edges {
		adj[e.From] = append(adj[e.From], e.To)
		indeg[e.To]++
	}
	var waves [][]Step
	processed := 0
	for processed < len(steps) {
		var wave []Step
		for id, deg := range indeg {
			if deg == 0 {
				wave = append(wave, steps[id])
				indeg[id] = -1
			}
		}
		if len(wave) == 0 {
			return nil, fmt.Errorf("cycle detected")
		}
		for _, s := range wave {
			processed++
			for _, next := range adj[s.ID] {
				indeg[next]--
			}
		}
		waves = append(waves, wave)
	}
	return waves, nil
}
