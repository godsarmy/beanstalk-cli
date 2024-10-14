package main

import (
    "context"
    "fmt"
    "net/url"

    "github.com/alecthomas/kingpin/v2"
    "github.com/beanstalkd/go-beanstalk"
)

var (
    debug = kingpin.Flag("debug", "Enable debug mode.").Short('d').Bool()
    address = kingpin.Flag("address", "Connect address.").Short('c').Default(":8080").String()

    bury = kingpin.Command("bury", "Bury a job.")
    buryJob = bury.Arg("job", "Job ID").Required().Uint64()

    delete = kingpin.Command("delete", "Delete a job.")
    deleteJob = delete.Arg("job", "Job ID").Required().Uint64()

    kick = kingpin.Command("kick", "Kick a job.")
    kickJob = kick.Arg("job", "Job ID").Required().Uint64()

    listt = kingpin.Command("list-tubes", "List Tubes.")

    peek = kingpin.Command("peek", "Peek a job.")
    peekJob = peek.Arg("job", "Job ID").Required().Uint64()

    peakReady = kingpin.Command("peek-ready", "Peek Ready Jobs.")
    peakReadyTube = peakReady.Arg("tube", "Tube name").String()

    put = kingpin.Command("put", "Put a job.")
    putBody = put.Arg("body", "Job body").Required().String()
    putPriority = put.Flag("priority", "Job priority").Short('p').Uint32()
    putTube = kingpin.Flag("tube", "Enable debug mode.").Short('t').String()
    putDelay = put.Flag("delay", "Job delay").Short('d').Duration()
    putTtr = put.Flag("ttr", "Job ttr").Short('r').Duration()

    release = kingpin.Command("release", "Release a job.")
    releaseJob = release.Arg("job", "Job ID").Required().Int()

    reserve = kingpin.Command("reserve", "Reserve a job.")
    reserveJob = reserve.Arg("job", "Job ID").Int()

    stats = kingpin.Command("stats", "Get stats.")

    statsj= kingpin.Command("stats-job", "Get job stats.")
    statsjJob = statsj.Arg("job", "Job ID").Required().Int()

    touch = kingpin.Command("touch", "Touch a job.")
    touchJob = touch.Arg("job", "Job ID").Required().Int()

)

func getConnect(address string) (*beanstalk.Conn, error) {
    u, err := url.Parse(address)
    if err != nil {
        return nil, err
    }

    if u.Scheme == "file" {
        return beanstalk.Dial("udp", u.Path)
    }

    if u.Scheme == "tcp" || u.Scheme == "" {
        return beanstalk.Dial("tcp", u.Host)
    }

    return nil, fmt.Errorf("Unknown scheme: %s", u.Scheme)
}

func reserveFunc(ctx context.Context, id uint64) {
    conn = ctx.Value("conn").(*beanstalk.Conn)
    body, err := conn.Reserve(*reserveJob)
    if err != nil {
        println(err)
        return
    }
    println(body)
}

func peekFunc(ctx context.Context, id uint64) {
    conn = ctx.Value("conn").(*beanstalk.Conn)
    body, err := conn.Peek(*peekJob)
    if err != nil {
        println(err)
        return
    }
    println(body)
}

func putFunc(ctx context.Context, body string, tube string, priority uint32, delay time.Duration, ttr time.Duration) {
    conn = ctx.Value("conn").(*beanstalk.Conn)
    id, err := conn.Put(tube, []byte(body), priority, delay, ttr)
    if err != nil {
        println(err)
        return
    }
    println(id)
}

func releaseFunc(ctx context.Context, id uint64) {
    conn = ctx.Value("conn").(*beanstalk.Conn)
    err := conn.Release(id)
    if err != nil {
        println(err)
        return
    }
}

func main() {
    ctx := context.Background()

    switch kingpin.Parse() {
        case "debug":
            ctx = context.WithValue(ctx, "debug", true)
        case "address":
            c, err := getConnect(*address)
            if err != nil {
                println(err)
                return
            }
            defer c.Close()
            ctx := context.WithValue(ctx, "conn", c)
       case "peek":
           peekFunc(ctx, *peekJob)
       case "reserve":
           reserveFunc(ctx, *reserveJob)
       case "put":
           putFunc(ctx, *putBody, *putTube, *putPriority, *putDelay, *putTtr)
       case "release":
           releaseFunc(ctx, *releaseJob)
       case "bury":
           buryFunc(ctx, *buryJob)
       case "delete":
           deleteFunc(ctx, *deleteJob)
       case "touch":
           touchFunc(ctx, *touchJob)
       case "stats":
           statsFunc(ctx, *statsjJob)
       case "stats-job":
           statsjFunc(ctx, *statsjJob)
    }
}
