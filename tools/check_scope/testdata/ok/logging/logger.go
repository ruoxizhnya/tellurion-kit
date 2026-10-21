// 夹具：一个干净的中性包。名字里没有任何 §2 词表的东西。
package logging

import (
	"io"
	"strings"
)

type Logger struct {
	w     io.Writer
	Level string
}

func New(w io.Writer) *Logger { return &Logger{w: w} }

func (l *Logger) WithLevel(level string) *Logger {
	l.Level = level
	return l
}

func MaskAPIKey(s string) string { return strings.ReplaceAll(s, "sk-", "***") }
