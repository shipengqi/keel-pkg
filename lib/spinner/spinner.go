package spinner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var characters = [4]string{"▖", "▘", "▝", "▗"}

type Spinner struct {
	*sync.RWMutex
	delay  time.Duration // delay is the speed of the indicator, default value 'time.Millisecond * 100'
	writer io.Writer     // writer default value os.Stdout
	prefix string
	suffix string
	endMsg string        // endMsg default value '\n'
	active bool          // active holds the state of the spinner
	stopC  chan struct{} // stopC is a channel used to stop the Spinner
}

// New create a  Spinner
func New() *Spinner {
	s := &Spinner{
		RWMutex: &sync.RWMutex{},
		writer:  os.Stdout,
		delay:   time.Millisecond * 100,
		endMsg:  "\n",
		active:  false,
		stopC:   make(chan struct{}, 0),
	}
	return s
}

func (s *Spinner) WithWriter(w io.Writer) *Spinner {
	s.writer = w
	return s
}

func (s *Spinner) WithDelay(delay time.Duration) *Spinner {
	s.delay = delay
	return s
}

func (s *Spinner) WithPrefix(prefix string) *Spinner {
	s.prefix = prefix
	return s
}

func (s *Spinner) WithSuffix(suffix string) *Spinner {
	s.suffix = suffix
	return s
}

func (s *Spinner) WithEndMsg(endMsg string) *Spinner {
	s.endMsg = endMsg
	return s
}

// Active whether spinner is currently active.
func (s *Spinner) Active() bool {
	return s.active
}

// Start will start the spinner.
// hides the cursor
// echo -e "\033[?25l"
// display the cursor
// echo -e "\033[?25h"
func (s *Spinner) Start() {
	s.Lock()
	if s.active {
		s.Unlock()
		return
	}
	s.active = true
	s.Unlock()

	go func() {
		for {
			for i := 0; i < len(characters); i++ {
				select {
				case <-s.stopC:
					return
				default:
					s.Lock()
					if !s.active {
						s.Unlock()
						return
					}
					plaintext := fmt.Sprintf("\r%s%s%s ", s.prefix, characters[i], s.suffix)
					_, _ = fmt.Fprint(s.writer, plaintext)
					s.Unlock()
					time.Sleep(s.delay)
				}
			}
		}
	}()
}

// Stop stops the spinner.
func (s *Spinner) Stop() {
	s.Lock()
	defer s.Unlock()
	if s.active {
		if len(s.endMsg) > 0 {
			_, _ = fmt.Fprint(s.writer, s.endMsg)
		}
		s.active = false
		s.stopC <- struct{}{}
		close(s.stopC)
	}
}

// Reset the spinner stop channel.
func (s *Spinner) Reset() *Spinner {
	s.stopC = make(chan struct{}, 0)
	s.active = false
	return s
}

// StopWithStatus stops the spinner and using the status instead of the spinner char.
func (s *Spinner) StopWithStatus(status string) string {
	s.Lock()
	defer s.Unlock()
	if s.active {
		plaintext := fmt.Sprintf("\r%s%s%s ", s.prefix, status, s.suffix)
		_, _ = fmt.Fprint(s.writer, plaintext)
		if len(s.endMsg) > 0 {
			_, _ = fmt.Fprint(s.writer, s.endMsg)
		}
		s.active = false
		s.stopC <- struct{}{}
		close(s.stopC)
		return plaintext
	}
	return ""
}
