package main

import (
	"fmt"
	"sync"
	"time"
)

type loader struct {
	isRunning bool
	speed     time.Duration
	design    []string
	preMsg    string
	posMsg    string
	endMsg    string
	mu        sync.Mutex // Mutex to synchronize access to isRunning
}

type Option func(*loader)

func WithDesign(spinnerdesign []string) Option {
	return func(l *loader) {
		l.design = spinnerdesign
	}
}
func WithSpeed(speedrate time.Duration) Option {
	return func(l *loader) {
		l.speed = speedrate
	}
}
func WithPreMsg(premsg string) Option {
	return func(l *loader) {
		l.preMsg = premsg
	}
}
func WithPosMsg(posmsg string) Option {
	return func(l *loader) {
		l.preMsg = posmsg
	}
}
func WithEndMsg(endmsg string) Option {
	return func(l *loader) {
		l.preMsg = endmsg
	}
}

func (l *loader) MainLoop() {
	i := 0
	for {
		l.mu.Lock()
		if !l.isRunning {
			l.mu.Unlock()
			break
		}
		l.mu.Unlock()

		fmt.Printf("\r%s%s%s ", l.preMsg, l.design[i], l.posMsg)
		i = (i + 1) % len(l.design)
		time.Sleep(l.speed)
	}
}

func (l *loader) Start() {
	l.mu.Lock()
	l.isRunning = true
	l.mu.Unlock()
	go l.MainLoop()
}

func (l *loader) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isRunning = false
	fmt.Printf("\r%s\n", l.endMsg)
}

func NewLoader(opts ...Option) *loader {
	l := loader{
		isRunning: false,
		design:    []string{"◜", "◝", "◞", "◟"},
		speed:     100 * time.Millisecond,
		preMsg:    "",
		posMsg:    "",
		endMsg:    "",
	}
	for _, opt := range opts {
		opt(&l)
	}
	return &l
}
