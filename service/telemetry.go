package service

import (
	"io"
	"strings"
	"time"
)

// LogWriter is a writer for disects logs before passing them along
type LogWriter struct {
	out  io.Writer
	tags map[string][]string
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	// Logging middleware

	// Split into tokens
	tokens := strings.Split(string(p), " ")

	// Tags are held in the 6th index
	tag := tokens[6][1 : len(tokens[6])-2] // strip parentheses
	if _, ok := lw.tags[tag]; ok == false {
		lw.tags[tag] = make([]string, 0)
	}
	lw.tags[tag] = append(lw.tags[tag], string(p))

	// Also categorize by worker id
	workerID := tokens[5]
	if _, ok := lw.tags[workerID]; ok == false {
		lw.tags[workerID] = make([]string, 0)
	}
	lw.tags[workerID] = append(lw.tags[workerID], string(p))

	// Write to output stream
	return lw.out.Write(p)
}

// Log in a standard format for service
func (s *Service) Log(tag string, message string) {
	s.Logger.Printf("%s %s (%s): %s", time.Now().String(), s.ID, tag, message)
}
