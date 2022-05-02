/*
 * @Descripttion:
 * @version:
 * @Author: wkq
 * @Date: 2022-05-01 11:33:19
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2022-05-02 21:55:21
 */
package main

import (
	mylog "go_study/consolelog"
)

const FILESIZE = 1024

func main() {
	// loger := mylog.NewLogger("debug")
	// loger.LogDebug("debug log")
	// loger.LogTrace("trace log")
	// loger.LogInfo("info log")
	// loger.LogWarning("warning log")
	// name := 121213
	// loger.LogError("error log name:%d", name)
	// loger.LogFatal("fatal log")
	filelog := mylog.NewFileLog("debug", "./", "filelog.log", "errfilelog", FILESIZE)
	for {
		filelog.LogDebug("debug log")
		filelog.LogTrace("trace log")
		filelog.LogInfo("info log")
		filelog.LogWarning("warning log")
		filelog.LogError("error log")
		filelog.LogFatal("fatal log")
	}

}
