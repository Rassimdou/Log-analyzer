package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogStorage struct {
	file      *os.File
	fileDate  string
	directory string
	mutex     sync.Mutex
}

func NewLogStorage(directory string) (*LogStorage, error) {
	// Create directory if it does not exist
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Get current date for file name (FIXED DATE FORMAT)
	currentDate := time.Now().Format("2006-01-02")
	filename := filepath.Join(directory, "security_"+currentDate+".log") // FIXED VARIABLE NAME

	// Open log file
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Return initialized LogStorage
	return &LogStorage{
		file:      file,
		fileDate:  currentDate,
		directory: directory,
	}, nil
}

func (s *LogStorage) WriteLog(method, source, ip, path, message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if we need to rotate logs (new day)
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != s.fileDate {
		if err := s.rotate(currentDate); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}
	}

	// Write log entry
	logEntry := formatLogEntry(method, source, ip, path, message)
	_, err := s.file.WriteString(logEntry) // Added comma after logEntry
	return err
}

// REMOVED THE NESTED FUNCTION DEFINITION
func (s *LogStorage) rotate(newDate string) error {
	// Close current file
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	// Create new file with new date
	filename := filepath.Join(s.directory, "security_"+newDate+".log")

	// Open new file
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	// Update storage with new file and date
	s.file = file
	s.fileDate = newDate
	return nil
}

func formatLogEntry(method, source, ip, path, message string) string {
	return fmt.Sprintf("[%s] [%s] [%s] %s\n",
		time.Now().UTC().Format(time.RFC3339),
		method,
		source,
		ip,
		path,
		message)
}
