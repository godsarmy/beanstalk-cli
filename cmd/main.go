package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/gosuri/uitable"
)

var version = "dev"

var (
	debug   = kingpin.Flag("debug", "Enable debug mode.").Short('d').Bool()
	address = kingpin.Flag("address", "Connect address.").
		Short('c').
		Default("tcp://127.0.0.1:11300").
		String()
	format = kingpin.Flag("format", "Output format").Default("text").Enum("json", "text")

	bury         = kingpin.Command("bury", "Bury a job.")
	buryJob      = bury.Arg("job", "Job ID").Required().Uint64()
	buryPriority = bury.Flag("priority", "Job priority").Short('p').Default("0").Uint32()

	delete    = kingpin.Command("delete", "Delete a job.")
	deleteJob = delete.Arg("job", "Job ID").Required().Uint64()

	kick    = kingpin.Command("kick", "Kick a job.")
	kickJob = kick.Arg("job", "Job ID").Required().Uint64()

	listt = kingpin.Command("list-tubes", "List Tubes.")

	peek    = kingpin.Command("peek", "Peek a job.")
	peekJob = peek.Arg("job", "Job ID").Required().Uint64()

	peakReady     = kingpin.Command("peek-ready", "Peek Ready Jobs.")
	peakReadyTube = peakReady.Arg("tube", "Tube name").String()

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

	touch    = kingpin.Command("touch", "Touch a job.")
	touchJob = touch.Arg("job", "Job ID").Required().Uint64()
)

func getConnect(address string) (*beanstalk.Conn, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "file" {
		return beanstalk.Dial("unix", u.Path)
	}

	if u.Scheme == "tcp" {
		return beanstalk.Dial("tcp", u.Host)
	}

	return nil, fmt.Errorf("Unknown scheme: %s", u.Scheme)
}

func printUI(result map[string]string) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	for key, value := range result {
		table.AddRow(fmt.Sprintf("%s:", key), value)
	}
	fmt.Println(table)
}

func printJSON(result map[string]string) {
	jsonString, err := json.Marshal(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(string(jsonString))
}

func reserveFunc(ctx context.Context, timeout time.Duration, tube string) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	if tube != "" {
		tubes := strings.Split(tube, ",")
		conn.TubeSet = *beanstalk.NewTubeSet(conn, tubes...)
	}
	id, body, err := conn.Reserve(timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]string{"id": strconv.FormatUint(id, 10), "body": string(body)}
}

func peekFunc(ctx context.Context, id uint64) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	body, err := conn.Peek(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]string{"body": string(body)}
}

func putFunc(
	ctx context.Context,
	body string,
	tube string,
	priority uint32,
	delay time.Duration,
	ttr time.Duration,
) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)

	if tube != "" {
		conn.Tube = *beanstalk.NewTube(conn, tube)
	}
	id, err := conn.Put([]byte(body), priority, delay, ttr)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]string{"id": strconv.FormatUint(id, 10)}
}

func releaseFunc(
	ctx context.Context,
	id uint64,
	priority uinte2,
	delay time.Duration,
) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	// job must be reserved first before release
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err := conn.Release(id, priority, delay)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]string{"id": strconv.FormatUint(id, 10), "body": string(body)}
}

func buryFunc(ctx context.Context, id uint64, priority uint32) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	// job must be reserved first before bury
	body, err := conn.ReserveJob(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = conn.Bury(id, priority)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]string{"id": strconv.FormatUint(id, 10), "body": string(body)}
}

func deleteFunc(ctx context.Context, id uint64) {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	err := conn.Delete(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func touchFunc(ctx context.Context, id uint64) {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	err := conn.Touch(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func statsFunc(ctx context.Context) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	stats, err := conn.Stats()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return stats
}

func statsjFunc(ctx context.Context, id uint64) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	stats, err := conn.StatsJob(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return stats
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

	var resp map[string]string
	switch cmd {
	case "peek":
		resp = peekFunc(ctx, *peekJob)
	case "reserve":
		resp = reserveFunc(ctx, *reserveTimeout, *reserveTube)
	case "put":
		resp = putFunc(ctx, *putBody, *putTube, *putPriority, *putDelay, *putTtr)
	case "release":
		resp = releaseFunc(ctx, *releaseJob, *releasePriority, *releaseDelay)
	case "bury":
		resp = buryFunc(ctx, *buryJob, *buryPriority)
	case "delete":
		deleteFunc(ctx, *deleteJob)
	case "touch":
		touchFunc(ctx, *touchJob)
	case "stats":
		resp = statsFunc(ctx)
	case "stats-job":
		resp = statsjFunc(ctx, *statsjJob)
	}

	if resp == nil {
		return
	}

	if *format == "json" {
		printJSON(resp)
	} else {
		printUI(resp)
	}
}
