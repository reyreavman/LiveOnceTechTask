package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     int                    `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service"`
	UserID    string                 `json:"user_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Logger - асинхронный логгер
type Logger struct {
	logChan    chan LogEntry
	done       chan struct{}
	wg         sync.WaitGroup
	minLevel   int
	service    string
	shutdown   bool
	shutdownMu sync.Mutex
}

// Создание нового логгера
func NewLogger() *Logger {
	logger := &Logger{
		logChan:  make(chan LogEntry, 1000),
		done:     make(chan struct{}),
		minLevel: INFO,
		service:  "my-service",
	}

	logger.wg.Add(1)
	go logger.processLogs()

	return logger
}

// Обработчик логов
func (l *Logger) processLogs() {
	defer l.wg.Done()

	for {
		select {
		case entry, ok := <-l.logChan:
			if !ok {
				return
			}
			l.writeLog(entry)
		case <-l.done:
			// Дописать оставшиеся логи перед выходом
			for {
				select {
				case entry := <-l.logChan:
					l.writeLog(entry)
				default:
					return
				}
			}
		}
	}
}

func (l *Logger) writeLog(entry LogEntry) {
	logData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Logger error: %v\n", err)
		return
	}
	os.Stdout.Write(append(logData, '\n'))
}

func (l *Logger) Debug(msg string, metadata map[string]interface{}) {
	l.log(DEBUG, msg, "", metadata)
}

func (l *Logger) Info(msg, userID string, metadata map[string]interface{}) {
	l.log(INFO, msg, userID, metadata)
}

func (l *Logger) Warning(msg, userID string, metadata map[string]interface{}) {
	l.log(WARNING, msg, userID, metadata)
}

func (l *Logger) Error(msg, userID string, err error, metadata map[string]interface{}) {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	metadata["error"] = err.Error()
	
	l.log(ERROR, msg, userID, nil)
}

func (l *Logger) log(level int, message, userID string, metadata map[string]interface{}) {
	l.shutdownMu.Lock()
	defer l.shutdownMu.Unlock()
	
	if l.shutdown || level < l.minLevel {
		return
	}

	if level == DEBUG {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			if metadata == nil {
				metadata = make(map[string]interface{})
			}
			metadata["caller"] = fmt.Sprintf("%s:%d", file, line)
		}
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   message,
		Service:   l.service,
		UserID:    userID,
		Metadata:  metadata,
	}

	select {
	case l.logChan <- entry:
	default:
		fmt.Fprintf(os.Stderr, "Logger buffer overflow. Message: %s\n", message)
	}
}

func (l *Logger) Stop() {
	l.shutdownMu.Lock()
	l.shutdown = true
	close(l.done)
	l.shutdownMu.Unlock()

	l.wg.Wait()
	close(l.logChan)
}