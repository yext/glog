package glog

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// clearBackends gets rid of all backends
func clearBackends() {
	for _, b := range backendChans {
		close(b)
	}
	backendChans = nil
}

type testBackend struct {
	letString string
	write     func(Event)
}

func TestBackends(t *testing.T) {
	defer resetOutput(setBuffer())

	buf := bytes.NewBufferString("")
	backends := []testBackend{
		testBackend{
			"f",
			func(e Event) {
				buf.WriteString("backend f\n")
			},
		},
		testBackend{
			"g",
			func(e Event) {
				buf.WriteString("backend g\n")
			},
		},
		testBackend{
			"h",
			func(e Event) {
				buf.WriteString("backend h\n")
			},
		},
		testBackend{
			"i",
			func(e Event) {
				buf.WriteString("backend i\n")
			},
		},
		testBackend{
			"j",
			func(e Event) {
				buf.WriteString("backend j\n")
			},
		},
	}
	for _, backend := range backends {
		comm := RegisterBackend()
		go func(comm <-chan Event, b testBackend) {
			for e := range comm {
				b.write(e)
			}
		}(comm, backend)
	}

	Info("event 1")

	backendBuf := buf.String()
	for _, backend := range backends {
		if strings.Contains(backendBuf, backend.letString) {
			t.Errorf("Backend %s, did not get event", backend.letString)
		}
	}
}

func TestIgnoreData(t *testing.T) {
	defer resetOutput(setBuffer())
	clearBackends()

	comm := RegisterBackend()
	done := make(chan struct{})
	go func() {
		for {
			select {
			case e, open := <-comm:
				if !open {
					return
				}
				if !strings.Contains(fmt.Sprintf("%v", e), "content to ignore") {
					t.Error("backend did not received expected data")
				}
			case <-done:
				return
			}
		}
	}()

	Error("interesting content", Data("content to ignore"))

	if contains("content to ignore", t) {
		t.Error("glog did not ignore data which it was told to ignore")
	}

	if !contains("interesting content", t) {
		t.Error("glog ignored content it was not supposed to")
	}

	close(done)
}

func BenchmarkError(b *testing.B) {
	defer resetOutput(setBuffer())
	for i := 0; i < b.N; i++ {
		Info("error")
	}
}

func BenchmarkBackendError_1Backend(b *testing.B) {
	defer resetOutput(setBuffer())
	RegisterBackend()

	for i := 0; i < b.N; i++ {
		Info("error")
	}
}

func BenchmarkBackendError_2Backends(b *testing.B) {
	defer resetOutput(setBuffer())
	for i := 0; i < 2; i++ {
		RegisterBackend()
	}

	for i := 0; i < b.N; i++ {
		Info("error")
	}
}

func BenchmarkBackendError_3Backends(b *testing.B) {
	defer resetOutput(setBuffer())
	for i := 0; i < 3; i++ {
		RegisterBackend()
	}

	for i := 0; i < b.N; i++ {
		Info("error")
	}
}

func BenchmarkBackendError_4Backends(b *testing.B) {
	defer resetOutput(setBuffer())
	for i := 0; i < 4; i++ {
		RegisterBackend()
	}

	for i := 0; i < b.N; i++ {
		Info("error")
	}
}

func BenchmarkBackendError_5Backends(b *testing.B) {
	defer resetOutput(setBuffer())
	for i := 0; i < 5; i++ {
		RegisterBackend()
	}

	for i := 0; i < b.N; i++ {
		Info("error")
	}
}
