package sod

import (
	"errors"
	"sync"
)

type runCount struct {
	sync.Mutex
	cnt int
}

var errMaxSolves = errors.New("max solvers have been run")

func (rc *runCount) grab() error {
	rc.Lock()
	defer rc.Unlock()
	if rc.cnt > 0 {
		rc.cnt--
		return nil
	}
	return errMaxSolves
}
