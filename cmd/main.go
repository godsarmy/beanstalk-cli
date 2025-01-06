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
	app     = kingpin.New("beanstalk-cli", "A beanstalkd client CLI application.")
	debug   = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	address = app.Flag("address", "Connect address. schema: tcp | unix").
		Envar("BS_ADDRESS").
		Short('a').
		Default("tcp://127.0.0.1:11300").
		String()
	format = app.Flag("format", "Output format").Short('f').Default("text").Enum("json", "text")

	bury         = app.Command("bury", "Bury a job.")
	buryJob      = bury.Arg("job", "Job ID").Required().Uint64()
	buryPriority = bury.Flag("priority", "Job priority").Short('p').Default("0").Uint32()

	delete    = app.Command("delete", "Delete a job.")
	deleteJob = delete.Arg("job", "Job ID").Required().Uint64()

	kick    = app.Command("kick", "Kick a job.")
	kickJob = kick.Arg("job", "Job ID").Required().Uint64()

	listt = app.Command("list-tubes", "List Tubes.")

	pause      = app.Command("pause-tube", "Pause new reservation for a Tube.")
	pauseTube  = pause.Arg("tube", "Tube name").String()
	pauseDelay = pause.Flag("delay", "Job delay").Required().Short('l').Duration()

	peek    = app.Command("peek", "Peek a job.")
	peekJob = peek.Arg("job", "Job ID").Required().Uint64()

	peekReady     = app.Command("peek-ready", "Peek Ready Jobs.")
	peekReadyTube = peekReady.Arg("tube", "Tube name").String()

	peekBuried     = app.Command("peek-buried", "Peek Buried Jobs.")
	peekBuriedTube = peekBuried.Arg("tube", "Tube name").String()

	peekDelayed     = app.Command("peek-deplayed", "Peek Delayed Jobs.")
	peekDelayedTube = peekDelayed.Arg("tube", "Tube name").String()

	put         = app.Command("put", "Put a job.")
	putBody     = put.Arg("body", "Job body").Required().String()
	putPriority = put.Flag("priority", "Job priority").Short('p').Uint32()
	putTube     = put.Flag("tube", "Enable debug mode.").Short('b').String()
	putDelay    = put.Flag("delay", "Job delay").Short('l').Duration()
	putTtr      = put.Flag("ttr", "Job ttr").Short('r').Duration()

	release         = app.Command("release", "Release a job.")
	releaseJob      = release.Arg("job", "Job ID").Required().Uint64()
	releasePriority = release.Flag("priority", "Job priority").Short('p').Default("0").Uint32()
	releaseDelay    = release.Flag("delay", "Job delay").Short('l').Default("0").Duration()

	reserve        = app.Command("reserve", "Reserve a job.")
	reserveTube    = reserve.Flag("tube", "Tube name").Short('b').String()
	reserveTimeout = reserve.Flag("timeout", "timeout").Short('t').Default("0").Duration()
	reserveJob     = reserve.Arg("job", "Job ID").Uint64()

	stats = app.Command("stats", "Get server stats.")

	statsj    = app.Command("stats-job", "Get job stats.")
	statsjJob = statsj.Arg("job", "Job ID").Required().Uint64()

	statst     = app.Command("stats-tube", "Get tube stats.")
	statstTube = statst.Arg("tube", "Tube name").Required().String()

	touch    = app.Command("touch", "Touch a job.")
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

	app.Version(version)
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
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
	case "list-tubes":
		resp := listTubesFunc(ctx)
		print(*format, resp)
	case "pause-tube":
		pauseTubeFunc(ctx, *pauseTube, *pauseDelay)
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
	case "put":
		resp := putFunc(ctx, *putBody, *putTube, *putPriority, *putDelay, *putTtr)
		print(*format, resp)
	case "release":
		resp := releaseFunc(ctx, *releaseJob, *releasePriority, *releaseDelay)
		print(*format, resp)
	case "reserve":
		var resp map[string]interface{}
		if *reserveJob != 0 {
			resp = reserveJobFunc(ctx, *reserveJob)
		} else {
			resp = reserveFunc(ctx, *reserveTimeout, *reserveTube)
		}
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
