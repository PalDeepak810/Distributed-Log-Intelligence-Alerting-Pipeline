package alert

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type ErrorTracker struct {
	mu         sync.Mutex
	counts     map[string]int
	timestamps map[string][]time.Time
}

var tracker = ErrorTracker{
	counts:     make(map[string]int),
	timestamps: make(map[string][]time.Time),
}

const (
	windowSize     = 10 * time.Second
	errorThreshold = 5
)

func CheckErrorSpike(service string) {

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	service = strings.TrimSpace(service)
	if service == "" {
		service = "unknown-service"
	}

	now := time.Now()

	// append current error timestamp
	tracker.timestamps[service] = append(tracker.timestamps[service], now)

	// remove old timestamps
	validTimes := []time.Time{}
	for _, t := range tracker.timestamps[service] {
		if now.Sub(t) <= windowSize {
			validTimes = append(validTimes, t)
		}
	}

	tracker.timestamps[service] = validTimes

	// check threshold
	if len(validTimes) >= errorThreshold {
		fmt.Printf(" SPIKE ALERT: %d errors in %v for service %s\n",
			len(validTimes),
			windowSize,
			service,
		)

		// reset after alert
		tracker.timestamps[service] = []time.Time{}
	}
}
