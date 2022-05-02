/*
 * @Descripttion:
 * @version:
 * @Author: wkq
 * @Date: 2022-05-01 22:53:15
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2022-05-02 21:59:54
 */
package consolelog

import (
	"fmt"
	"os"
	"path"
	"time"
)

type Filelog struct {
	level       LogLevel
	filepath    string //日志存放文件路径
	filename    string //日志存放文件名
	errfilename string //错误日志存放文件名
	fileobj     *os.File
	errfileobj  *os.File
	maxSize     int64
}

//构造函数
func NewFileLog(levelstr, fp, fn, errfn string, size int64) *Filelog {
	ll, err := ParseLogLevel(levelstr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	f1 := &Filelog{
		level:       ll,
		filepath:    fp,
		filename:    fn,
		errfilename: errfn,
		maxSize:     size,
	}
	fileerr := f1.initFile()

	if fileerr != nil {
		fmt.Println(fileerr)
		panic(fileerr)
	}
	return f1
}

//检查文件大小
func (f *Filelog) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return fileInfo.Size() > f.maxSize
}

//打开文件
func (f *Filelog) initFile() error {
	url := path.Join(f.filepath, f.filename)
	filedata, err1 := os.OpenFile(url, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		return err1
	}
	errfiledata, err2 := os.OpenFile(url+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		return err2
	}
	f.errfileobj = errfiledata
	f.fileobj = filedata
	return nil
}

//关闭文件
func (f *Filelog) CloseFile() {
	defer f.fileobj.Close()
	defer f.errfileobj.Close()
}

//控制日志层级
func (f *Filelog) isAble(level LogLevel) bool {
	return level >= f.level
}

//切割文件
func (f *Filelog) splitFile(file *os.File) (*os.File, error) {
	nowstr := time.Now().Format("20060102150405000")
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	logName := path.Join(f.filepath, fileInfo.Name())               //当前日志文件路径
	new_log_filename := fmt.Sprintf("%s.chunk.%s", logName, nowstr) //拼接日志文件备份的name
	// 1.关闭当前 日志文件
	file.Close()
	//2.拷贝文件 xx.log --->xx.log.bak1223

	os.Rename(logName, new_log_filename)
	//3.打开一个新的日志文件
	filedata, err2 := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err2 != nil {
		return nil, err2
	}
	//4.将新的日志文件对象赋值给f.fileobj
	return filedata, nil
}

//日志写到文件中
func (f *Filelog) log(l LogLevel, formate string, arg ...interface{}) {
	if f.isAble(l) {
		now := time.Now()
		msg := fmt.Sprintf(formate, arg...)
		funcName, filepath, line := GetRowNum(3)
		timestr := now.Format("2006-01-02 15:04:05")
		if f.checkSize(f.fileobj) { //判断是否需要切割
			newfile, err := f.splitFile(f.fileobj)
			if err != nil {
				return
			}
			f.fileobj = newfile
		}

		fmt.Fprintf(f.fileobj, "time:[%s] - type:%s - filepath:%s - funcName:%s - line:%d - info:%s\n", timestr, levleMap[uint16(l)], filepath, funcName, line, msg)
		if l >= ERROR {
			if f.checkSize(f.errfileobj) { //判断是否需要切割
				newfile, err := f.splitFile(f.errfileobj)
				if err != nil {
					return
				}
				f.errfileobj = newfile
			}
			fmt.Fprintf(f.errfileobj, "time:[%s] - type:%s - filepath:%s - funcName:%s - line:%d - info:%s\n", timestr, levleMap[uint16(l)], filepath, funcName, line, msg)
		}
		// f.CloseFile()
	}
}

func (f Filelog) LogDebug(formate string, arg ...interface{}) {
	f.log(DEBUG, formate, arg...)
}

//trace
func (f Filelog) LogTrace(formate string, arg ...interface{}) {
	f.log(TRACE, formate, arg...)
}

//info
func (f Filelog) LogInfo(formate string, arg ...interface{}) {
	f.log(INFO, formate, arg...)
}

//warning
func (f Filelog) LogWarning(formate string, arg ...interface{}) {
	f.log(WARNING, formate, arg...)
}

//error
func (f Filelog) LogError(formate string, arg ...interface{}) {
	f.log(ERROR, formate, arg...)
}

//fatal
func (f Filelog) LogFatal(formate string, arg ...interface{}) {
	f.log(FATAL, formate, arg...)
}
