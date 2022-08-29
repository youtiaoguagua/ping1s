package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"time"
)

const dbUrl = "https://ghproxy.com/https://raw.githubusercontent.com/youtiaoguagua/ping1s/master/ping1s.db"

func initDB() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(
				color.Error,
				"[ %v ] %s\n",
				color.New(color.FgRed, color.Bold).Sprint("ERROR"),
				err,
			)
			os.Exit(int(exitCodeErrDB))
		}
	}()

	exists := dbExists(homeDir)
	if !exists {
		startDownload()
	} else {
		// 判断数据库是否损坏
	}

}

func startDownload() {
	s := spinner.New(spinner.CharSets[14], 500*time.Millisecond)
	s.Color("blue")
	s.Prefix = color.New(color.FgHiYellow, color.Bold).Sprintf("初始化数据: ")
	s.Start()
	downloadDb(homeDir)
	s.Stop()
}

func downloadDb(homeDir string) {
	// download db
	res, err := http.Get(dbUrl)
	if err != nil {
		log.Error(err)
		return
	}
	defer res.Body.Close()

	os.Chmod(path.Join(homeDir, "/.ping1s"), fs.ModePerm)

	out, err := os.Create(path.Join(homeDir, "/.ping1s", dbName))
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer out.Close()

	// save db
	_, err = io.Copy(out, res.Body)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func dbExists(homeDir string) bool {
	_, err := os.Stat(path.Join(homeDir, "/.ping1s", dbName))
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
