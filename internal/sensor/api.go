package sensor

import (
	"github.com/torbenschinke/i2c"
	"github.com/torbenschinke/picd/pkg/logging"
	"sync"
	"time"
)

type SenseService struct {
	lastT  i2c.T
	lastRH i2c.RH
	poller *i2c.Polling
	mutex  sync.Mutex
}

func NewSenseService() *SenseService {
	s := &SenseService{
		lastT:  0,
		lastRH: 0,
		poller: i2c.NewPolling(10 * time.Second),
	}

	s.poller.Dispatcher().Register(s)

	return s
}

func (s *SenseService) T() i2c.T {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastT
}

func (s *SenseService) RH() i2c.RH {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.lastRH
}

func (s *SenseService) OnTemperature(id i2c.ID, t i2c.T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastT = t
}

func (s *SenseService) OnError(id i2c.ID, err error) {
	logging.FromContext(nil).Println(id, err)
}

func (s *SenseService) OnHumidity(id i2c.ID, h i2c.RH) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastRH = h
}
