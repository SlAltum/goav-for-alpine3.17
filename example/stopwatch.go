package main

import (
	"time"
)

type Stopwatch struct {
	start time.Time
	stop  time.Time
}

var stopwatch *Stopwatch

func GetStopwatch() *Stopwatch {
	if stopwatch == nil {
		stopwatch = new(Stopwatch)
	}
	return stopwatch
}

func (s *Stopwatch) Start() {
	s.start = time.Now()
}

func (s *Stopwatch) Stop() {
	s.stop = time.Now()
}

//纳秒
func (s Stopwatch) RuntimeNs() int64 {

	return s.stop.UnixNano() - s.start.UnixNano()
}

//微妙
func (s Stopwatch) RuntimeUs() float64 {
	return (float64)(s.RuntimeNs()) / 1000.00
}

//毫秒
func (s Stopwatch) RuntimeMs() float64 {
	return (float64)(s.RuntimeNs()) / 1000000.00
}

//秒
func (s Stopwatch) RuntimeS() float64 {
	return (float64)(s.RuntimeNs()) / 1000000000.00
}
