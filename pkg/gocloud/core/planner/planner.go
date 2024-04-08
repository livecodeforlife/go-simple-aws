package planner

import "github.com/livecodeforlife/go-simple-aws/pkg/gocloud/core"

// NewSimplePlanner creates a new planner to be used inside the infra package
func NewSimplePlanner() core.Planner {
	return &simplePlanner{}
}

type simplePlanner struct {
	resources []core.LazyResourceInterface
}

func (p *simplePlanner) AddResource(resource core.LazyResourceInterface) error {
	p.resources = append(p.resources, resource)
	return nil
}

func (p *simplePlanner) TopoSortForCreation() []core.LazyResourceInterface {
	return p.resources
}

func (p *simplePlanner) TopoSortForDeletion() []core.LazyResourceInterface {
	return reverseSliceCopy(p.resources)
}

func reverseSliceCopy[T any](s []T) []T {
	sc := make([]T, len(s))
	copy(sc, s)
	for i, j := 0, len(sc)-1; i < j; i, j = i+1, j-1 {
		sc[i], sc[j] = sc[j], sc[i]
	}
	return sc
}
