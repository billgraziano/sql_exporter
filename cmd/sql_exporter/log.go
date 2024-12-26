package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/common/promslog"
)

type logConfig struct {
	logger         *slog.Logger
	logFileHandler *os.File
}

// initLogFile opens the log file for writing if a log file is specified.
func initLogFile(logFile string) (*os.File, error) {
	if logFile == "" {
		return nil, nil
	}
	logFile = logFileWithTimeStamp(logFile)
	logFileHandler, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}
	return logFileHandler, nil
}

// logFileWithTimeStamp inserts a YYYY_MM timestamp in a file name
func logFileWithTimeStamp(logFile string) string {
	ext := filepath.Ext(logFile)
	name := strings.TrimSuffix(logFile, ext)
	tsname := fmt.Sprintf("%s_%d_%d%s", name, time.Now().Year(), time.Now().Month(), ext)
	return tsname
}

// initLogConfig configures and initializes the logging system.
func initLogConfig(logLevel, logFormat string, logFile string) (*logConfig, error) {
	logFileHandler, err := initLogFile(logFile)
	if err != nil {
		return nil, err
	}

	if logFileHandler == nil {
		logFileHandler = os.Stderr
	}

	promslogConfig := &promslog.Config{
		Level:  &promslog.AllowedLevel{},
		Format: &promslog.AllowedFormat{},
		Style:  promslog.SlogStyle,
		Writer: logFileHandler,
	}

	if err := promslogConfig.Level.Set(logLevel); err != nil {
		return nil, err
	}

	if err := promslogConfig.Format.Set(logFormat); err != nil {
		return nil, err
	}
	// Initialize logger.
	logger := promslog.New(promslogConfig)

	return &logConfig{
		logger:         logger,
		logFileHandler: logFileHandler,
	}, nil
}
