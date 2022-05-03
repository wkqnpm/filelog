
package consolelog

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

type LogLevel uint16

const (
	UNNOW LogLevel = iota
	ALL
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

var levleMap = map[uint16]string{2: "debug", 3: "trace", 4: "info", 5: "warning", 6: "error", 7: "fatal"}

//Logger
type Logger struct {
	level LogLevel
}

//NewLogger
func NewLogger(level string) Logger {
	ll, err := ParseLogLevel(level)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return Logger{
		level: ll,
	}
}

//控制日志层级
func (l Logger) isAble(level LogLevel) bool {
	return level >= l.level
}

//解析日志级别
func ParseLogLevel(s string) (LogLevel, error) {
	switch s {
	case "all":
		return ALL, nil
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("invalid log level")
		return UNNOW, err
	}
}

//生成日志x
func (c Logger) log(l LogLevel, formate string, arg ...interface{}) {
	if c.isAble(l) {
		now := time.Now()
		msg := fmt.Sprintf(formate, arg...)
		funcName, filepath, line := GetRowNum(3)
		timestr := now.Format("2006-01-02 15:04:05")
		fmt.Printf("time:[%s] - type:%s - filepath:%s - funcName:%s - line:%d - info:%s\n", timestr, levleMap[uint16(l)], filepath, funcName, line, msg)
	}
}

//获取日志行号
func GetRowNum(n int) (funcName, filepath string, line int) {
	/*
		n 为0 line为当前行
		n 为1 向上找一层
		n 为2 向上找两层
	*/

	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		fmt.Println("runtime.Caller() faild\n")
		return
	}
	fN := runtime.FuncForPC(pc).Name()
	return fN, file, line
}

//debug
func (l Logger) LogDebug(formate string, arg ...interface{}) {
	l.log(DEBUG, formate, arg...)

}

//trace
func (l Logger) LogTrace(msg string) {
	l.log(TRACE, msg)
}

//info
func (l Logger) LogInfo(msg string) {
	l.log(INFO, msg)
}

//warning
func (l Logger) LogWarning(msg string) {
	l.log(WARNING, msg)
}

//error
func (l Logger) LogError(formate string, arg ...interface{}) {
	l.log(ERROR, formate, arg...)
}

//fatal
func (l Logger) LogFatal(msg string) {
	l.log(FATAL, msg)
}
