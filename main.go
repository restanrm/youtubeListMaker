package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

func createYoutubeLink(src io.Reader) {
	var ids []string
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		url, err := url.Parse(line)
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err, "line": line}).Error("This line is not a proper url")
			continue
		}
		values := url.Query()
		if id, ok := values["v"]; !ok {
			logrus.WithField("line", line).Error("This line is not properly formed. Could not find id of video")
		} else {
			ids = append(ids, id[0])
		}
	}
	if len(ids) <= 0 {
		logrus.Error("Could not parse a single line to produce a youtube playlist")
		os.Exit(-1)
	} else {
		var playListUrl = "https://www.youtube.com/watch_videos?video_ids=" + strings.Join(ids, ",")
		fmt.Println("URL of the playlist: ", playListUrl)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() > 0 {
		if _, err := os.Stat(flag.Arg(0)); err != nil {
			logrus.WithField("err", err).Error("Error while opening youtube link list")
			os.Exit(-1)
		}
		fh, err := os.Open(flag.Arg(0))
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err, "file": flag.Arg(0)}).Error("Could not open file")
			os.Exit(-1)
		}
		createYoutubeLink(fh)
	} else {
		createYoutubeLink(os.Stdin)
	}
}
