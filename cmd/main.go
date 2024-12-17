package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/gosuri/uitable"
)

var version = "dev"

var (
	debug   = kingpin.Flag("debug", "Enable debug mode.").Short('d').Bool()
	address = kingpin.Flag("address", "Connect address. schema: tcp | unix").
		Short('a').
		Default("tcp://127.0.0.1:11300").
		String()
	format = kingpin.Flag("format", "Output format").Short('f').Default("text").Enum("json", "text")

	bury         = kingpin.Command("bury", "Bury a job.")
	buryJob      = bury.Arg("job", "Job ID").Required().Uint64()
	buryPriority = bury.Flag("priority", "Job priority").Short('p').Default("0").Uint32()

	delete    = kingpin.Command("delete", "Delete a job.")
	deleteJob = delete.Arg("job", "Job ID").Required().Uint64()

	kick    = kingpin.Command("kick", "Kick a job.")
	kickJob = kick.Arg("job", "Job ID").Required().Uint64()

	listt = kingpin.Command("list-tubes", "List Tubes.")

	pause      = kingpin.Command("pause", "Pause new reservation for a Tube.")
	pauseTube  = pause.Arg("tube", "Tube name").String()
	pauseDelay = pause.Flag("delay", "Job delay").Required().Short('l').Duration()

	peek    = kingpin.Command("peek", "Peek a job.")
	peekJob = peek.Arg("job", "Job ID").Required().Uint64()

	peekReady     = kingpin.Command("peek-ready", "Peek Ready Jobs.")
	peekReadyTube = peekReady.Arg("tube", "Tube name").String()

	peekBuried     = kingpin.Command("peek-buried", "Peek Buried Jobs.")
	peekBuriedTube = peekBuried.Arg("tube", "Tube name").String()

	peekDelayed     = kingpin.Command("peek-deplayed", "Peek Delayed Jobs.")
	peekDelayedTube = peekDelayed.Arg("tube", "Tube name").String()

	put         = kingpin.Command("put", "Put a job.")
	putBody     = put.Arg("body", "Job body").Required().String()
	putPriority = put.Flag("priority", "Job priority").Short('p').Uint32()
	putTube     = put.Flag("tube", "Enable debug mode.").Short('b').String()
	putDelay    = put.Flag("delay", "Job delay").Short('l').Duration()
	putTtr      = put.Flag("ttr", "Job ttr").Short('r').Duration()

	release         = kingpin.Command("release", "Release a job.")
	releaseJob      = release.Arg("job", "Job ID").Required().Uint64()
	releasePriority = release.Flag("priority", "Job priority").Short('p').Default("0").Uint32()
	releaseDelay    = release.Flag("delay", "Job delay").Short('l').Default("0").Duration()

	reserve        = kingpin.Command("reserve", "Reserve a job.")
	reserveTube    = reserve.Flag("tube", "Tube name").Short('b').String()
	reserveTimeout = reserve.Flag("timeout", "timeout").Short('t').Default("0").Duration()

	stats = kingpin.Command("stats", "Get stats.")

	statsj    = kingpin.Command("stats-job", "Get job stats.")
	statsjJob = statsj.Arg("job", "Job ID").Required().Uint64()

	statst     = kingpin.Command("stats-tube", "Get tube stats.")
	statstTube = statst.Arg("tube", "Tube name").Required().String()

	touch    = kingpin.Command("touch", "Touch a job.")
	touchJob = touch.Arg("job", "Job ID").Required().Uint64()
)

func getConnect(address string) (*beanstalk.Conn, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "unix" {
		return beanstalk.Dial("unix", u.Path)
	}

	if u.Scheme == "tcp" {
		return beanstalk.Dial("tcp", u.Host)
	}

	return nil, fmt.Errorf("Unknown scheme: %s", u.Scheme)
}

func printUI[T any](result map[string]T) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	for key, value := range result {
		table.AddRow(fmt.Sprintf("%s:", key), value)
	}
	fmt.Println(table)
}

func printJSON[T any](result map[string]T) {
	jsonString, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(string(jsonString))
}

func print[T any](format string, data map[string]T) {
	if data == nil {
		return
	}

	if format == "json" {
		printJSON(data)
	} else {
		printUI(data)
	}
}

func main() {
	ctx := context.Background()

	kingpin.Version(version)
	cmd := kingpin.Parse()
	ctx = context.WithValue(ctx, "debug", *debug)
	c, err := getConnect(*address)
	if err != nil {
		fmt.Println(err)
		os.Exit(255)
	}
	defer c.Close()
	ctx = context.WithValue(ctx, "conn", c)

	switch cmd {
	case "bury":
		resp := buryFunc(ctx, *buryJob, *buryPriority)
		print(*format, resp)
	case "delete":
		deleteFunc(ctx, *deleteJob)
	case "kick":
		kickFunc(ctx, *kickJob)
	case "peek":
		resp := peekFunc(ctx, *peekJob)
		print(*format, resp)
	case "peekBuried":
		resp := peekBuriedFunc(ctx, *peekBuriedTube)
		print(*format, resp)
	case "peekDelayed":
		resp := peekDelayedFunc(ctx, *peekDelayedTube)
		print(*format, resp)
	case "peekReady":
		resp := peekReadyFunc(ctx, *peekReadyTube)
		print(*format, resp)
	case "pause":
		pauseFunc(ctx, *pauseTube, *pauseDelay)
	case "put":
		resp := putFunc(ctx, *putBody, *putTube, *putPriority, *putDelay, *putTtr)
		print(*format, resp)
	case "release":
		resp := releaseFunc(ctx, *releaseJob, *releasePriority, *releaseDelay)
		print(*format, resp)
	case "reserve":
		resp := reserveFunc(ctx, *reserveTimeout, *reserveTube)
		print(*format, resp)
	case "stats":
		resp := statsFunc(ctx)
		print(*format, resp)
	case "stats-job":
		resp := statsjFunc(ctx, *statsjJob)
		print(*format, resp)
	case "stats-tube":
		resp := statstFunc(ctx, *statstTube)
		print(*format, resp)
	case "touch":
		touchFunc(ctx, *touchJob)
	}
}
