package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/mattn/go-runewidth"
	"strings"
)

func handlerPing1s(index int) string {

	if len(*poetryList) <= index {
		return strings.Repeat(" ", runewidth.StringWidth((*poetryList)[0]))
	}

	res := (*poetryList)[index]

	return res
}

func pingRecv(pkt *ping.Packet) {
	fmt.Printf("seq=%s %sbytes from %s: ttl=%s time=%s %s\n",
		color.New(color.FgHiYellow, color.Bold).Sprintf("%d", pkt.Seq),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
		color.New(color.FgHiCyan, color.Bold).Sprintf("%d", pkt.Ttl),
		color.New(color.FgHiMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
		handlerPing1s(pkt.Seq),
	)
}

func pingFinish(stats *ping.Statistics) {
	color.New(color.FgWhite, color.Bold).Printf(
		"\n───────── %s ping statistics ─────────\n",
		stats.Addr,
	)
	fmt.Printf(
		"%s: %v transmitted => %v received (%v loss)\n",
		color.New(color.FgHiWhite, color.Bold).Sprintf("PACKET STATISTICS"),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%d", stats.PacketsSent),
		color.New(color.FgHiGreen, color.Bold).Sprintf("%d", stats.PacketsRecv),
		color.New(color.FgHiRed, color.Bold).Sprintf("%v%%", stats.PacketLoss),
	)
	fmt.Printf(
		"%s: min=%v avg=%v max=%v stddev=%v\n",
		color.New(color.FgHiWhite, color.Bold).Sprintf("ROUND TRIP"),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%v", stats.MinRtt),
		color.New(color.FgHiCyan, color.Bold).Sprintf("%v", stats.AvgRtt),
		color.New(color.FgHiGreen, color.Bold).Sprintf("%v", stats.MaxRtt),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", stats.StdDevRtt),
	)
	hitokotoContent := (*hitokotoList)[0]
	fmt.Printf("\n%s %v %s \n",
		color.New(color.FgHiYellow, color.Bold).Sprintf("『"),
		color.New(color.FgHiWhite, color.Bold).Sprintf(hitokotoContent),
		color.New(color.FgHiYellow, color.Bold).Sprintf("』"),
	)

	from := fmt.Sprintf("—— 「%s」", (*hitokotoList)[1])
	count := runewidth.StringWidth(hitokotoContent+strings.Repeat(" ", 8)) + -runewidth.StringWidth(from)
	from = strings.Repeat(" ", count) + from
	fmt.Printf("%v \n", color.New(color.FgHiBlue, color.Bold).Sprintf(from))

	fmt.Println()

}
