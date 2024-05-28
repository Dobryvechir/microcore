/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvlog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var streamFile *os.File
var streamFileSync = false
var streamCrLf = []byte{byte(13), byte(10)}
var streamFileLock = &sync.Mutex{}

func CheckDvLogStreamToFile(parameters map[string]string) {
	if streamFile != nil {
		return
	}
	var err error
	fileName := parameters["MICROCORE_DVLOG_STREAM"]
	if fileName != "" {
		os.RemoveAll(fileName)
		streamFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			streamFile = nil
		}
	}
}

func GetCurrentTimeStampForLog() []byte {
	res := make([]byte, 9)
	now := time.Now()
	hour := now.Hour()
	minutes := now.Minute()
	seconds := now.Second()
	res[0] = byte((hour / 10) + '0')
	res[1] = byte((hour % 10) + '0')
	res[2] = ':'
	res[3] = byte((minutes / 10) + '0')
	res[4] = byte((minutes % 10) + '0')
	res[5] = ':'
	res[6] = byte((seconds / 10) + '0')
	res[7] = byte((seconds % 10) + '0')
	res[8] = ' '
	return res
}

func PrintlnToStream(s string) (err error) {
	if streamFile != nil {
		data := []byte(s)
		n := len(data)
		if n > 0 {
			timeStamp := GetCurrentTimeStampForLog()
			streamFileLock.Lock()
			streamFile.Write(timeStamp)
			_, err = streamFile.Write(data)
			if data[n-1] != 13 && data[n-1] != 10 {
				streamFile.Write(streamCrLf)
			}
			if !streamFileSync {
				streamFileSync = true
				go func() {
					time.Sleep(time.Duration(7) * time.Second)
					streamFileLock.Lock()
					streamFile.Sync()
					streamFileSync = false
					streamFileLock.Unlock()
				}()
			}
			streamFileLock.Unlock()
		}
	}
	return
}

func FlushStream() {
	if streamFile != nil {
		streamFileLock.Lock()
		streamFile.Sync()
		streamFileLock.Unlock()
	}
}

func Println(short string, full string) {
	if streamFile == nil {
		log.Println(full)
	} else {
		err := PrintlnToStream(full)
		if err == nil {
			if short != "" {
				log.Println(short)
			}
		} else {
			log.Println(full)
		}

	}
}

func PrintlnError(s string) {
	Println(s, s)
}

func Print(s string) {
	Println(s, s)
}


func PrintError(err error) {
   if err!=nil {
        PrintlnError(err.Error())
   }
}

func Printf(formatShort, formatFull string, a ...interface{}) {
	short := fmt.Sprintf(formatShort, a...)
	full := fmt.Sprintf(formatFull, a...)
	Println(short, full)
}

func PrintfFullOnly(format string, a ...interface{}) {
	full := fmt.Sprintf(format, a...)
	Println("", full)
}

func PrintfError(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	Println(s, s)
}
