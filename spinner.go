package spinner

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Spinner struct {
	mu         *sync.Mutex
	charSet    []string
	pos        int
	active     bool
	stopChan   chan struct{}
	suffix     string
	prefix     string
	lastOutput string
	delay      time.Duration
}

func (s *Spinner) next() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.charSet[s.pos%len(s.charSet)]
	s.pos++
	return r
}

func (s *Spinner) WithSuffix(suffix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.suffix = suffix
}

func (s *Spinner) WithPrefix(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prefix = prefix
}

func New(cs []string) *Spinner {
	return &Spinner{
		charSet:  cs,
		pos:      0,
		active:   false,
		stopChan: make(chan struct{}, 1),
		mu:       &sync.Mutex{},
		delay:    100 * time.Millisecond,
	}
}

func (s *Spinner) UpdateCharSet(cs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.charSet = cs
}

func (s *Spinner) SetDelay(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.delay = d
}

func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		return
	}
	s.active = true

	go func() {
		for {
			select {
			case <-s.stopChan:
				return
			default:
				fmt.Printf("\r%s%s%s ", s.prefix, s.next(), s.suffix)
				s.lastOutput = fmt.Sprintf("\r%s%s%s ", s.prefix, s.next(), s.suffix)
				time.Sleep(s.delay)
				s.erase()
			}
		}
	}()
}

func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	clearString := "\r" + strings.Repeat(" ", n) + "\r"
	fmt.Print(clearString)
	s.lastOutput = ""
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		s.erase()
		s.stopChan <- struct{}{}
	}
}
