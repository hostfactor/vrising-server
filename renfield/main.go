package main

import (
	"context"
	"fmt"
	"github.com/nxadm/tail"
	"os"
	"time"
)

const (
	ServerLogPrefix   = "[vrising-server] "
	RenfieldLogPrefix = "[renfield] "
)

var LogPath = os.Getenv("LOG_FILE")

func main() {
	fmt.Printf("%sStarting!\n", RenfieldLogPrefix)
	ctx := context.Background()

	go ReadyChecker(ctx)
	go TailFile(ctx, LogPath)
	<-ctx.Done()
	fmt.Printf("%sExiting.\n", RenfieldLogPrefix)
}

func ReadyChecker(ctx context.Context) {
	tickerDur := 10 * time.Second
	ticker := time.NewTicker(tickerDur)
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			fmt.Println(RenfieldLogPrefix + "Exiting ready checker")
			return
		case <-ticker.C:
			_, err := os.Stat(LogPath)
			if err != nil {
				fmt.Printf("%sServer not ready. This typically takes around ~5m. Elapsed time: %s.\n", RenfieldLogPrefix, time.Since(start))
			} else {
				fmt.Printf("%sYour V Rising server is ready! Took %s to start.\n", RenfieldLogPrefix, time.Since(start))
				return
			}
		}
	}
}

func TailFile(ctx context.Context, fp string) {
	tailer, err := tail.TailFile(fp, tail.Config{
		Follow: true,
	})
	if err != nil {
		panic(fmt.Errorf("failed to tail file %s: %w", fp, err))
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println(RenfieldLogPrefix + "Exiting file watcher.")
			return
		case line := <-tailer.Lines:
			fmt.Println(ServerLogPrefix + line.Text)
		}
	}
}
