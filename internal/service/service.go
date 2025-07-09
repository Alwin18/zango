package service

import (
	"fmt"
	"sync"
)

type Service struct {
	Name   string
	Active bool
}

var (
	services = map[string]*Service{
		"http": {Name: "http", Active: false},
		"db":   {Name: "db", Active: false},
	}
	mu sync.Mutex
)

func Start(name string) error {
	mu.Lock()
	defer mu.Unlock()

	svc, ok := services[name]
	if !ok {
		return fmt.Errorf("service %s tidak ditemukan", name)
	}

	if svc.Active {
		return fmt.Errorf("service %s sedang aktif", name)
	}

	svc.Active = true
	return nil
}

func Stop(name string) error {
	mu.Lock()
	defer mu.Unlock()

	svc, ok := services[name]
	if !ok {
		return fmt.Errorf("Service %s tidak ditemukan", name)
	}

	if !svc.Active {
		return fmt.Errorf("Service %s sudah berhenti", name)
	}

	svc.Active = false
	return nil
}

func Status(name string) (bool, error) {
	mu.Lock()
	defer mu.Unlock()

	svc, ok := services[name]
	if !ok {
		return false, fmt.Errorf("Service %s tidak ditemukan", name)
	}
	return svc.Active, nil
}
