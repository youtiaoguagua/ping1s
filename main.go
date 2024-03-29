package main

import (
	"fmt"
	"github.com/fatih/color"
	probing "github.com/prometheus-community/pro-bing"
	flag "github.com/spf13/pflag"
	"io/fs"
	"os"
	"os/signal"
	"path"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrPing
	exitCodeErrDB
	exitCodeErrUnexpected
)

const (
	dbName  = "ping1s.db"
	version = "v1.3"
)

var (
	homeDir     string
	commandArgs CommandArgs
)

var (
	poetryList, hitokotoList *[]string
)

func main() {
	initDB()
	code, err := start()
	if err != nil {
		stdOutErr(err)
	}

	os.Exit(int(code))

}

func stdOutErr(err error) {
	_, _ = fmt.Fprintf(
		color.Error,
		"[ %v ] %s\n",
		color.New(color.FgRed, color.Bold).Sprint("ERROR"),
		err,
	)
}

func start() (exitCode, error) {
	flag.StringVarP(&commandArgs.Type, "type", "t", "-1", `hitokoto type
a-动画 b-漫画 c-游戏 d-文学 e-原创 f-来自网络 
g-其他 h-影视 i-诗词 j-网易云 k-哲学 l-抖机灵
`)
	flag.StringVarP(&commandArgs.Author, "author", "a", "-1", "poetry author, such as 苏东坡")
	flag.IntVarP(&commandArgs.CollectionType, "collection", "c", -1, `collection of poetry
 0-诗经 1-唐诗 2-宋词 
3-教科书 4-花间集 5-南唐二主词
`)
	flag.IntVarP(&commandArgs.Num, "num", "n", 10, "number of poetry")
	flag.BoolVarP(&commandArgs.Version, "version", "v", false, "ping1s version")
	flag.Usage = usage
	flag.Parse()

	if commandArgs.Version {
		fmt.Printf("ping1s version: %s\n", version)
		return exitCodeOK, nil
	}

	if flag.NArg() == 0 {
		return exitCodeErrArgs, fmt.Errorf("host arg is required")
	} else if flag.NArg() > 1 {
		return exitCodeErrArgs, fmt.Errorf("too many args,there should be only one")
	}

	// cover collection type
	coverCollection()
	// query db
	poetryList, hitokotoList = startQueryDb()

	host := flag.Arg(0)
	pinger, err := initPing(host)

	if err != nil {
		return exitCodeOK, fmt.Errorf("an error occurred while initializing ping1s: %w", err)
	}

	if err := pinger.Run(); err != nil {
		return exitCodeErrPing, fmt.Errorf("an error occurred when running ping1s: %w", err)
	}

	return exitCodeOK, nil
}

func coverCollection() {
	switch commandArgs.CollectionType {
	case 0:
		commandArgs.Collection = "诗经"
	case 1:
		commandArgs.Collection = "唐诗"
	case 2:
		commandArgs.Collection = "宋词"
	case 3:
		commandArgs.Collection = "教科书"
	case 4:
		commandArgs.Collection = "花间集"
	case 5:
		commandArgs.Collection = "南唐二主词"
	default:
		commandArgs.Collection = ""
	}
}

func initPing(host string) (*probing.Pinger, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed to init pinger %w", err)
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()

	_, _ = color.New(color.FgHiWhite, color.Bold).Printf(
		"PING %s (%s) type `Ctrl-C` to abort\n",
		pinger.Addr(),
		pinger.IPAddr(),
	)

	// receivedPacket is a callback function that will
	pinger.OnRecv = pingRecv

	pinger.OnFinish = pingFinish

	return pinger, nil
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `Usage: ping1s [Options] HOST

Options:
`)
	flag.PrintDefaults()
}

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		stdOutErr(err)
		os.Exit(int(exitCodeErrUnexpected))
	}

	err = os.MkdirAll(path.Join(homeDir, "/.ping1s"), fs.ModePerm)

	if err != nil {
		stdOutErr(err)
		os.Exit(int(exitCodeErrUnexpected))
	}
	initLog()
}
