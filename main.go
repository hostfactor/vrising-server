package main

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/mmcdole/gofeed"
	"os"
	"os/exec"
	"regexp"
)

var semVer = regexp.MustCompile("\\d+\\.\\d+\\.\\d+")

var URLs = []string{
	"hostfactor/vrising-server",
	"ghcr.io/hostfactor/vrising-server",
}

func main() {
	fp := gofeed.NewParser()

	feed, _ := fp.ParseURL("https://store.steampowered.com/feeds/news/app/1604030/?cc=US&l=english&snr=1_2108_9__2107")

	var sv string
	var link string
	for _, item := range feed.Items {
		if sv != "" {
			break
		}

		sv = semVer.FindString(item.Title)
		if sv == "" {
			continue
		}

		_, err := semver.NewVersion(sv)
		if err != nil {
			sv = ""
			continue
		}
		link = item.Link
	}

	opts := []string{
		"build", "-t",
		"--build-arg", fmt.Sprintf("%s=%s", "VERSION", sv),
		"--build-arg", fmt.Sprintf("%s=%s", "VERSION_URL", link),
	}

	for _, t := range URLs {
		opts = append(opts, t+":"+sv)
		opts = append(opts, t+":latest")
	}

	opts = append(opts, ".")

	cmd := exec.Command("docker", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
