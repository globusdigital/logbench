package zap

import (
	"io"
	"net/url"
	"sync"

	"go.uber.org/zap"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	writer     io.ReadWriter
	lock       sync.Mutex
	registered bool
}

// Close implements the Sink interface
func (s *MemorySink) Close() error { return nil }

// Sync implements the Sink interface
func (s *MemorySink) Sync() error { return nil }

// Write implements the Sink interface
func (s *MemorySink) Write(data []byte) (int, error) {
	zapSink.lock.Lock()
	defer zapSink.lock.Unlock()
	return s.writer.Write(data)
}

// Read implements the io.ReadWriter interface
func (s *MemorySink) Read(data []byte) (int, error) {
	zapSink.lock.Lock()
	defer zapSink.lock.Unlock()
	return s.writer.Read(data)
}

// SetOut atomically sets the output read-writer
func (s *MemorySink) SetOut(out io.ReadWriter) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.writer = out
	if !s.registered {
		s.registered = true
		return zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
			return s, nil
		})
	}

	return nil
}

var _ io.ReadWriter = new(MemorySink)
var zapSink = new(MemorySink)
