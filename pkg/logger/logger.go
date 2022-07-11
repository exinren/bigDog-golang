package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

// 定义日志等级
type Level int8
// 日志的字段
type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// 判断当前等级。
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type Logger struct {
	newLogger *log.Logger
	ctx context.Context
	fields Fields
	callers []string
}

// Logger构造函数。
func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

// 复制结构体变量地址。
func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

//WithLevel：设置日志等级。

//WithFields：设置日志公共字段。
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if nil == ll.fields {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

//WithContext：设置日志上下文属性。
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}
//WithCaller：设置当前某一层调用栈的信息（程序计数器、文件信息、行号）。
func (l *Logger) WithCaller(skin int) *Logger {
	ll := l.clone()
	// Caller报告当前go程调用栈所执行的函数的文件和行号信息
	pc, file, line, ok := runtime.Caller(skin)
	if ok {
		//FuncForPC返回一个表示调用栈标识符pc对应的调用栈的*Func
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s：%d %s", file,line,f.Name())}
	}
	return ll
}
//WithCallersFrames：设置当前的整个调用栈信息。
func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	//实参skip为开始在pc中记录之前所要跳过的栈帧数，1表示Callers所在的调用栈。
	depth := runtime.Callers(minCallerDepth, pcs)
	// 获取 Callers 返回的 PC 值的一部分，并准备返回函数/文件/行信息。 在完成帧之前不要更改切片。
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callers = append(callers,fmt.Sprintf("%s：%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}

// 日志格式化和输出
func (l *Logger) JSONFormat(level Level, message string) map[string]interface{}  {
	data := make(Fields, len(l.fields) + 4)
	data["level"] = level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

// 根据等级和信息，输出响应的日志打印格式。
func (l *Logger) Output(level Level, message string) {
	body, _ := json.Marshal(l.JSONFormat(level, message))
	content := string(body)
	switch level {
	case LevelDebug:
		fallthrough
	case LevelInfo:
		fallthrough
	case LevelWarn:
		fallthrough
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

func (l *Logger) Debug(v...interface{}) {
	l.Output(LevelDebug, fmt.Sprint(v...))
}


func (l *Logger) Debugf(format string, v...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}


func (l *Logger) Infof(format string, v...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(v...interface{}) {
	l.Output(LevelWarn, fmt.Sprint(v...))
}


func (l *Logger) Warnf(format string, v...interface{}) {
	l.Output(LevelWarn, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v...interface{}) {
	l.Output(LevelError, fmt.Sprint(v...))
}


func (l *Logger) Errorf(format string, v...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}


func (l *Logger) Fatal(v...interface{})  {
	l.Output(LevelFatal, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
}


func (l *Logger) Panic(v...interface{})  {
	l.Output(LevelPanic, fmt.Sprint(v...))
}

func (l *Logger) Panicf(format string, v...interface{}) {
	l.Output(LevelPanic, fmt.Sprintf(format, v...))
}

