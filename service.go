package async_flag

import (
	"sync"
)

type Flag interface {
	TrySet() bool
	IsSet() bool
	Drop()
}

type service struct {
	mu   *sync.Mutex
	flag int
}

func New() Flag {
	mu := sync.Mutex{}
	s := service{
		mu:   &mu,
		flag: 0,
	}
	return &s
}

func (s *service) TrySet() bool {
	s.mu.Lock()
	if s.flag != 0 {
		s.mu.Unlock()
		return false
	}
	s.flag = 1
	s.mu.Unlock()
	return true
}

func (s *service) IsSet() bool {
	s.mu.Lock()
	if s.flag == 0 {
		s.mu.Unlock()
		return false
	}
	s.mu.Unlock()
	return true
}

func (s *service) Drop() {
	s.mu.Lock()
	s.flag = 0
	s.mu.Unlock()
}

func (s *service) TryDrop() bool {
	s.mu.Lock()
	if s.flag != 0 {
		s.flag = 0
		s.mu.Unlock()
		return true
	}
	s.mu.Unlock()
	return false
}
