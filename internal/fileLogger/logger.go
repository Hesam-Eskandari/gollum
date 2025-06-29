package fileLogger

import (
	"fmt"
	"log/slog"
	"path"
	"strings"
	"time"
)

const (
	maxCheckPeriod  = time.Hour
	maxPurgePeriod  = time.Hour * 24 * 30
	minCheckPeriod  = time.Second
	minPurgePeriod  = time.Second * 10
	maxFileSizeByte = 1024 * 1024 * 1024
	minFileSizeByte = 1024 * 100
	timeFormat      = "2006-01-02T15-04-05-999Z07-00"
)

type Config struct {
	Filename        string        // filename prefix
	PurgePeriod     time.Duration // max time wait to purge
	CheckPeriod     time.Duration // period to check the file size
	MaxFileSizeByte int64         // maximum allowed filesize before rotate
	MinLogLevel     slog.Level    // lowest allowed log level
}

func BuildDefaultConfig(filename string) Config {
	return Config{
		Filename:        filename,
		PurgePeriod:     time.Hour * 24 * 7,
		CheckPeriod:     time.Minute,
		MaxFileSizeByte: 1024 * 1024 * 64,
		MinLogLevel:     slog.LevelInfo,
	}
}

type Logger interface {
	// Setup configures a slog logger and rotates log file based on a config
	Setup() error
	// Destroy stops the logging process and frees the memory
	Destroy() error
}

func New(config Config) Logger {
	return &logger{
		config:      config,
		filename:    config.Filename,
		minLogLevel: config.MinLogLevel,
	}
}

type logger struct {
	config      Config
	filename    string
	minLogLevel slog.Level
	rw          *rotateWriter
}

// Setup configures a slog logger and rotates log file based on a config
func (l *logger) Setup() error {
	var err error
	l.rw, err = newRotateWriter(l.filename, l.config)
	if err != nil {
		return err
	}
	handler := slog.NewJSONHandler(
		l.rw,
		&slog.HandlerOptions{AddSource: false})
	slogger := slog.New(handler)
	slog.SetLogLoggerLevel(l.minLogLevel)
	slog.SetDefault(slogger)
	go l.rw.Check()
	return nil
}

// Destroy stops the logging process and frees the memory
func (l *logger) Destroy() error {
	return l.rw.destroy()
}

func getFilepath(filename string) string {
	dir, filename := path.Split(filename)
	extension := path.Ext(filename)
	fileBase := strings.Trim(filename, extension)
	return fmt.Sprintf("%v%v_%v%v", dir, fileBase, time.Now().Format(timeFormat), extension)
}
