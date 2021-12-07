package linuxcnc

import (
	"sync"
	"time"
)

type StatusPoller struct {
	c        statusPoller
	interval time.Duration

	onChangeCh    chan struct{}
	lastStatus    *Status
	currentStatus *Status
	statusMux     sync.RWMutex
}

type statusPoller interface {
	PollStatus(s *Status) error
}

func NewStatusPoller(c statusPoller, interval time.Duration) *StatusPoller {
	return &StatusPoller{
		c:             c,
		lastStatus:    &Status{},
		currentStatus: &Status{},
		interval:      interval,

		onChangeCh: make(chan struct{}),
	}
}

// Informs a client about changed status from the poller.
// Clients may call .Status() to get new status information.
func (s *StatusPoller) OnChange() <-chan struct{} {
	return s.onChangeCh
}

// Status returns a copy of the current status.
func (s *StatusPoller) Status() Status {
	s.statusMux.RLock()
	defer s.statusMux.RUnlock()
	return *s.currentStatus
}

func (s *StatusPoller) Run(stopCh <-chan struct{}) error {
	t := time.NewTicker(s.interval)
	defer t.Stop()

	if err := s.run(); err != nil {
		return err
	}

	for range t.C {
		if err := s.run(); err != nil {
			return err
		}
	}
	return nil
}

func (s *StatusPoller) run() error {
	s.statusMux.Lock()
	defer s.statusMux.Unlock()

	current := s.lastStatus
	s.lastStatus = s.currentStatus
	s.currentStatus = current

	if err := s.c.PollStatus(s.currentStatus); err != nil {
		return err
	}

	if !s.currentStatus.Equal(s.lastStatus) {
		s.onChangeCh <- struct{}{}
	}

	return nil
}
