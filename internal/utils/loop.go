package utils

import (
	"sync"
	"time"
)

func Loop(backoff time.Duration, fn func()) func() {
	cancelled := false
	var mu sync.Mutex

	cancel := func() {
		mu.Lock()
		cancelled = true
		mu.Unlock()
	}

	isCancelled := func() bool {
		mu.Lock()
		defer mu.Unlock()
		return cancelled
	}

	for !isCancelled() {
		fn()
		time.Sleep(backoff)
	}

	return cancel
}
