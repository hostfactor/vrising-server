package main

import (
	"context"
	"fmt"
	"github.com/nxadm/tail"
	"net/http"
	"os"
	"time"
)

const (
	ServerLogPrefix   = "[vrising-server] "
	RenfieldLogPrefix = "[renfield] "
)

var LogPath = os.Getenv("LOG_FILE")

var HttpPort = os.Getenv("RENFIELD_SERVER_PORT")

var IsReady = false

func main() {
	fmt.Printf("%sStarting!\n", RenfieldLogPrefix)
	ctx := context.Background()

	go ReadyChecker(ctx)
	go TailFile(ctx, LogPath)
	if HttpPort != "" {
		go HttpServer(ctx)
	}
	<-ctx.Done()
	fmt.Printf("%sExiting.\n", RenfieldLogPrefix)
}

func HttpServer(_ context.Context) {
	http.HandleFunc("/api/server/ready", func(writer http.ResponseWriter, request *http.Request) {
		if IsReady {
			writer.WriteHeader(http.StatusOK)
		} else {
			writer.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	http.HandleFunc("/api/server/backup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%s", HttpPort)
	fmt.Printf("%sStarting server on %s.\n", RenfieldLogPrefix, addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("%sServer failed to start: %s.\n", RenfieldLogPrefix, err.Error())
	}
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
				IsReady = true
				fmt.Printf("%sYour V Rising server is ready! Took %s to start.\n", RenfieldLogPrefix, time.Since(start))
				return
			}
		}
	}
}

func TailFile(ctx context.Context, fp string) {
	if fp == "" {
		return
	}

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
