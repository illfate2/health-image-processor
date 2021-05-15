package body

import (
	"go.uber.org/atomic"

	"context"
	"log"
	"time"
)

type Processor struct {
	lastFaceNotRecognizedNotified time.Time
	isFaceRecognized              *atomic.Bool
	isNoseRecognized              *atomic.Bool
	isBackRecognized              *atomic.Bool
	faceIsNotRecognizedCh         chan struct{}

	userLastBlinked   time.Time
	lastBlinkNotified time.Time
	blinkedTime       int
	blinkTimeout      int

	noseCrookedCh           chan struct{}
	lastCrookedNoseNotified time.Time
	noseCrookedTime         int
	noseCrookedTimeout      int

	backCrookedCh           chan struct{}
	lastCrookedBackNotified time.Time
	backCrookedTime         int
	backCrookedTimeout      int
}

func NewService() *Processor {
	return &Processor{
		isFaceRecognized:      atomic.NewBool(true),
		isBackRecognized:      atomic.NewBool(true),
		isNoseRecognized:      atomic.NewBool(true),
		faceIsNotRecognizedCh: make(chan struct{}),
		blinkedTime:           5,
		blinkTimeout:          5,
		noseCrookedCh:         make(chan struct{}, 100000),
		backCrookedCh:         make(chan struct{}, 100000),
	}
}

type NotifyMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func newNotifyMessage(t string, message string) NotifyMessage {
	return NotifyMessage{Type: t, Message: message}
}

func timeSinceMoreThenSec(from time.Time, than int) bool {
	return time.Since(from).Seconds() > float64(than)
}

const faceIsNotRecognizedTimeout = 10

func (s *Processor) StartNotifying(ctx context.Context) <-chan NotifyMessage {
	const ticketDur = time.Millisecond * 50
	blinkTicker := time.NewTicker(ticketDur)

	ch := make(chan NotifyMessage)
	go func() {
		for {
			select {
			case <-blinkTicker.C:
				if timeSinceMoreThenSec(s.userLastBlinked, s.blinkedTime) && timeSinceMoreThenSec(s.lastBlinkNotified, s.blinkTimeout) && s.isFaceRecognized.Load() {
					log.Print("Sending notification to blink")
					ch <- NotifyMessage{
						Type:    "blink",
						Message: "Pls, blink",
					}
					s.lastBlinkNotified = time.Now()
				}

			case <-s.backCrookedCh:
				if timeSinceMoreThenSec(s.lastCrookedBackNotified, 5) && s.isBackRecognized.Load() {
					ch <- NotifyMessage{
						Type:    "crookedBack",
						Message: "Pls, sit straight",
					}
					s.lastCrookedBackNotified = time.Now()
				}

			case <-s.noseCrookedCh:
				if timeSinceMoreThenSec(s.lastCrookedNoseNotified, 5) && s.isNoseRecognized.Load() {
					ch <- NotifyMessage{
						Type:    "crookedHead",
						Message: "Pls, set your head straight",
					}
					s.lastCrookedNoseNotified = time.Now()
				}
			case <-ctx.Done():
				blinkTicker.Stop()
				close(ch)
				return
			}
		}
	}()
	return ch
}

func (s *Processor) Blinked(isFaceRecognized bool) {
	if !isFaceRecognized {
		s.isFaceRecognized.Store(false)
		log.Print("face is not recognized for blinking")
		return
	}
	s.isFaceRecognized.Store(true)
	log.Print("Blinking...")
	s.userLastBlinked = time.Now()
}

func (s *Processor) NoseCrooked(isCrooked, isFaceRecognized bool) {
	if !isFaceRecognized {
		s.isNoseRecognized.Store(false)
		log.Print("face is not recognized for nose")
		return
	}
	s.isNoseRecognized.Store(true)
	log.Print("Nose crooked...", isCrooked)
	if isCrooked {
		s.noseCrookedCh <- struct{}{}
	}
}

func (s *Processor) BackCrooked(isCrooked, isFaceRecognized bool) {
	if !isFaceRecognized {
		s.isBackRecognized.Store(false)
		log.Print("face is not recognized for back")
		return
	}
	s.isBackRecognized.Store(true)
	log.Print("Back crooked...", isCrooked)
	if isCrooked {
		s.backCrookedCh <- struct{}{}
	}
}
