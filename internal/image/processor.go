package image

import "time"

type Processor struct {
	ch         chan struct{}
	timeoutSec int
	lastTime   time.Time
}

func NewProcessor() *Processor {
	return &Processor{
		ch:         make(chan struct{}),
		timeoutSec: 1,
		lastTime:   time.Now(),
	}
}

func (p *Processor) GetNotificationCh() <-chan struct{} {
	return p.ch
}

func (p *Processor) ProcessImage(matrix [][]int) {
	if isUserSittingCrooked(matrix) && time.Since(p.lastTime).Seconds() > float64(p.timeoutSec) {
		p.ch <- struct{}{}
		p.lastTime = time.Now()
	}
}

func isUserSittingCrooked(matrix [][]int) bool {
	return false
}
