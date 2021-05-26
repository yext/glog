package glog

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

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

	comm := RegisterBackend()

	message := fmt.Sprintf("testIgnoreData message: %v", time.Now().Nanosecond())
	Error(message, Data("data1"))

	if contains("data1", t) {
		t.Error("glog did not ignore data which it was told to ignore")
	}

	if !contains(message, t) {
		t.Error("glog ignored content it was not supposed to")
	}

	waitForData(t, comm, message, "data1")
}

func TestFormatString(t *testing.T) {
	defer resetOutput(setBuffer())

	comm := RegisterBackend()

	message := "error: test error"
	formatMsg := "error: %s"
	Errorf(formatMsg, errors.New("test error"))

	if !contains(message, t) {
		t.Error("glog did not process errorf to stdout")
	}
	waitForData(t, comm, message, FormatStringArg{formatMsg})
}

func TestErrorArgs(t *testing.T) {
	defer resetOutput(setBuffer())

	comm := RegisterBackend()

	err := errors.New("test error")
	Error(err)

	if !contains(err.Error(), t) {
		t.Error("glog did not process error message to stdout")
	}
	waitForData(t, comm, err.Error(), ErrorArg{err})
}

func waitForData(t *testing.T, comm <-chan Event, expectedMessage string, expectedData ...interface{}) {
	timeout := time.After(1 * time.Second)
	for {
		select {
		case e, open := <-comm:
			if !open {
				t.Error("Backend closed without receiving expected message")
				return
			}
			if strings.Contains(string(e.Message), expectedMessage) {
				set := make(map[interface{}]struct{})
				for _, d := range e.Data {
					t.Log(d)
					set[d] = struct{}{}
				}
				for _, d := range expectedData {
					if _, ok := set[d]; !ok {
						t.Error("Expected data not received in message: ", d)
					}
				}
				return
			}
		case <-timeout:
			t.Error("Timed out waiting for data on backend")
			return
		}
	}
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
