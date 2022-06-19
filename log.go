package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func initLog() {
	file := path.Join(homeDir, ".ping1s/ping1s.log")
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	stat, err := logFile.Stat()
	if err != nil {
		panic(err)
	}
	fileWriter := Logger{logFile, stat.Size()}
	log.SetOutput(&fileWriter)
	log.SetOutput(logFile)
	log.SetLevel(log.WarnLevel)
	log.SetReportCaller(true)
}

func (l *Logger) Write(data []byte) (n int, err error) {
	file := path.Join(homeDir, ".ping1s/ping1s.log")
	if l == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if l.file == nil {
		return 0, errors.New("file not opened")
	}
	n, e := l.file.Write(data)
	l.size += int64(n)

	if l.size > 1*1024*1024 {
		err = l.file.Close()
		if err != nil {
			return 0, err
		}

		err = os.Remove(file)
		if err != nil {
			return 0, err
		}

		l.file, _ = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			return 0, err
		}
		l.size = 0
	}
	return n, e
}
