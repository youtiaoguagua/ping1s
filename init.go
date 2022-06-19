package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const dbUrl = "https://wxfwxfwxf.coding.net/p/ping1sping1s/d/ping1s/git/raw/master/ping1s.db"

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
		s := spinner.New(spinner.CharSets[14], 500*time.Millisecond)
		s.Color("blue")
		s.Prefix = color.New(color.FgHiYellow, color.Bold).Sprintf("init ping1s: ")
		s.Start()
		downloadDb(homeDir)
		s.Stop()
	}

}

func downloadDb(homeDir string) {
	// download db
	res, err := http.Get(dbUrl)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	err = os.MkdirAll(path.Join(homeDir, "/.ping1s"), 777)

	if err != nil {
		panic(err)
	}

	out, err := os.Create(path.Join(homeDir, "/.ping1s", dbName))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// save db
	_, err = io.Copy(out, res.Body)
	if err != nil {
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
