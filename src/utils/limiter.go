package utils

import (
	log "github.com/cihub/seelog"
	"sync"
)

type Limiter struct {
	maxInflight     int
	currentInflight int
	mutex           *sync.Mutex
}

func (l *Limiter) Get() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.maxInflight == 0 {
		return true
	} else if l.currentInflight >= l.maxInflight {
		return false
	} else {
		l.currentInflight += 1
		return true
	}
}

func (l *Limiter) Release() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.maxInflight == 0 || l.currentInflight <= 0 {
		return
	} else {
		l.currentInflight -= 1
		return
	}
}

func NewLimiter(maxInflight int) *Limiter {
	limit := &Limiter{
		maxInflight:     maxInflight,
		currentInflight: 0,
		mutex:           new(sync.Mutex),
	}
	log.Infof("Initialized Limiter with maxInflight = %v.", limit.maxInflight)
	return limit
}
