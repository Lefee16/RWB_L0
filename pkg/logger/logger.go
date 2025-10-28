package logger

import (
	"log"
)

// Logger - интерфейс логера
type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// SimpleLogger - простая реализация логера
type SimpleLogger struct {
	level string
}

// New - создание нового логера
func New(level string) Logger {
	return &SimpleLogger{level: level}
}

// Info - информационное сообщение
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if l.shouldLog("info") {
		log.Printf("[INFO] "+msg, args...)
	}
}

// Warn - предупреждение
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if l.shouldLog("warn") {
		log.Printf("[WARN] "+msg, args...)
	}
}

// Error - ошибка
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if l.shouldLog("error") {
		log.Printf("[ERROR] "+msg, args...)
	}
}

// Debug - отладочное сообщение
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if l.shouldLog("debug") {
		log.Printf("[DEBUG] "+msg, args...)
	}
}

// shouldLog - проверка, нужно ли логировать на текущем уровне
func (l *SimpleLogger) shouldLog(msgLevel string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	currentLevel := levels[l.level]
	messageLevel := levels[msgLevel]

	return messageLevel >= currentLevel
}
