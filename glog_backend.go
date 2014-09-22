package glog

import "runtime"

var (
	messageChan  = make(chan Event, 10)
	backendChans []chan<- Event
)

type data struct {
	d interface{}
}

// Data tags an item to be ignored by glog when logging but to pass it
// to any registered backends. An example would be:
//
//   var req *http.Request
//   ...
//   glog.Error("failed to complete process", glog.Data(req))
func Data(arg interface{}) interface{} {
	return data{arg}
}

// An Event contains a logged event's severity (INFO, WARN, ERROR, FATAL),
// a format string (if Infof, Warnf, Errorf or Fatalf were used) and a slice
// of everything else passed to the log call.
type Event struct {
	Severity   string
	Message    []byte
	Data       []interface{}
	StackTrace []uintptr
}

// NewEvent creates a glog.Event from the logged event's severity,
// format string (if Infof, Warnf, Errorf or Fatalf were called) and
// any other arguments passed to the log call.
// NewEvent separates out any items tagged by Data() and stores them
// in Event.Data.
func NewEvent(s severity, message []byte, dataArgs []interface{}, extraDepth int, frames []uintptr) Event {
	var stackTrace []uintptr

	if s >= errorLog {
		callers := make([]uintptr, 20)
		written := runtime.Callers(4+extraDepth, callers)
		stackTrace = callers[:written]
	}

	stackTrace = append(frames, stackTrace...)
	return Event{
		Severity:   severityName[s],
		Message:    message,
		Data:       dataArgs,
		StackTrace: stackTrace,
	}
}

// filterData splits out any items tagged by Data() and returns two slices:
// the first with only argments meant for the log call and the second with
// only arguments meant to passed to any registered backends.
func filterData(args []interface{}) ([]interface{}, []interface{}) {
	var (
		realArgs []interface{}
		dataArgs []interface{}
	)

	for _, arg := range args {
		if argd, ok := arg.(data); ok {
			dataArgs = append(dataArgs, argd.d)
		} else {
			realArgs = append(realArgs, arg)
		}
	}
	return realArgs, dataArgs
}

// RegisterBackend returns a channel on which Event's will be passed
// when they are logged.
//
// The caller is responsible for any necessary synchronization such
// that the call to this function "happens before" any events to be
// logged to this channel or other calls to RegisterBackend().
func RegisterBackend() <-chan Event {
	if len(backendChans) == 0 {
		go broadcastEvents()
	}

	c := make(chan Event, 100)
	backendChans = append(backendChans, c)
	return c
}

// eventForBackends creates and writes a glog.Event to the message channel
// if and only if we have registered backends.
func eventForBackends(e Event) {
	if len(backendChans) > 0 {
		select {
		case messageChan <- e:
		default:
		}
	}
}

func broadcastEvents() {
	for e := range messageChan {
		for _, c := range backendChans {
			select {
			case c <- e:
			default:
			}
		}
	}
}
