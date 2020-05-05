package services

import (
	"time"

	"github.com/sirupsen/logrus"
)

type services struct {
	logger *logrus.Logger
}

// NewServices will create new an services object to start bamboo services.
func NewServices(logger *logrus.Logger) *services {
	return &services{
		logger: logger,
	}
}

// Run trying to execute function f 'repeat' of times with time interval 'timestamp'
func (s services) Run(repeat int, timestamp int, f func() error) {
	for repeat > 0 {
		if err := f(); err == nil {
			return
		}
		repeat--

		time.Sleep(time.Duration(timestamp) * time.Second)
	}
}

// Forever trying to execute function f with time interval 'timestamp'
func (s services) Forever(timestamp int, f func() error) {
	for {
		if err := f(); err != nil {
			return
		}

		time.Sleep(time.Duration(timestamp) * time.Second)
	}
}
