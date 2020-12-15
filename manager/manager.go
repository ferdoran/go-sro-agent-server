package manager

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Manager struct {
	ticker       *time.Ticker
	startTimer   *time.Timer
	started      bool
	mutex        sync.Mutex
	Name         string
	runnerFunc   func()
	initialDelay time.Duration
	rate         time.Duration
}

func (m *Manager) Start() {
	logrus.Infof("starting %s in %v and a rate of %v", m.Name, m.initialDelay, m.rate)
	if m.ticker != nil && m.started {
		logrus.Warnf("already started KnownListManager")
		return
	}

	if m.startTimer != nil {
		logrus.Warnf("already starting %s", m.Name)
		return
	}

	m.startTimer = time.AfterFunc(m.initialDelay, func() {
		m.ticker = time.NewTicker(m.rate)
		m.mutex.Lock()
		defer m.mutex.Unlock()
		m.started = true
		go m.runnerFunc()
		logrus.Infof("started %s", m.Name)
	})
}

func (m *Manager) Stop() {
	logrus.Infof("stopping %s", m.Name)
	if m.startTimer != nil {
		m.mutex.Lock()
		m.startTimer.Stop()
		m.startTimer = nil
		m.mutex.Unlock()
	}

	if m.ticker == nil {
		logrus.Warnf("%s has already been stopped", m.Name)
		m.mutex.Lock()
		defer m.mutex.Unlock()
		m.started = false
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.started = false
	m.ticker.Stop()
	m.ticker = nil
}
