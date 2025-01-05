package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

func reserveJobFunc(ctx context.Context, id uint64) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	body, err := conn.ReserveJob(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"body": string(body)}
}

func reserveFunc(ctx context.Context, timeout time.Duration, tube string) map[string]interface{} {
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
	return map[string]interface{}{"id": id, "body": string(body)}
}

func peekFunc(ctx context.Context, id uint64) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	body, err := conn.Peek(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"body": string(body)}
}

func putFunc(
	ctx context.Context,
	body string,
	tube string,
	priority uint32,
	delay time.Duration,
	ttr time.Duration,
) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)

	if tube != "" {
		conn.Tube = *beanstalk.NewTube(conn, tube)
	}
	id, err := conn.Put([]byte(body), priority, delay, ttr)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"id": id}
}

func releaseFunc(
	ctx context.Context,
	id uint64,
	priority uint32,
	delay time.Duration,
) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	// job must be reserved first before release
	body, err := conn.ReserveJob(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = conn.Release(id, priority, delay)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"id": id, "body": string(body)}
}

func buryFunc(ctx context.Context, id uint64, priority uint32) map[string]interface{} {
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
	return map[string]interface{}{"id": id, "body": string(body)}
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

func statstFunc(ctx context.Context, tube string) map[string]string {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	connTube := *beanstalk.NewTube(conn, tube)
	stats, err := connTube.Stats()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return stats
}

func kickFunc(ctx context.Context, id uint64) {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	err := conn.KickJob(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

func pauseTubeFunc(ctx context.Context, tube string, delay time.Duration) {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	connTube := *beanstalk.NewTube(conn, tube)
	err := connTube.Pause(delay)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

func peekReadyFunc(ctx context.Context, tube string) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	if tube != "" {
		conn.Tube = *beanstalk.NewTube(conn, tube)
	}

	id, body, err := conn.PeekReady()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"id": id, "body": string(body)}
}

func peekDelayedFunc(ctx context.Context, tube string) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	if tube != "" {
		conn.Tube = *beanstalk.NewTube(conn, tube)
	}

	id, body, err := conn.PeekDelayed()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"id": id, "body": string(body)}
}

func peekBuriedFunc(ctx context.Context, tube string) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	if tube != "" {
		conn.Tube = *beanstalk.NewTube(conn, tube)
	}

	id, body, err := conn.PeekBuried()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"id": id, "body": string(body)}
}

func listTubesFunc(ctx context.Context) map[string]interface{} {
	conn := ctx.Value("conn").(*beanstalk.Conn)
	tubes, err := conn.ListTubes()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return map[string]interface{}{"tubes": tubes}
}
