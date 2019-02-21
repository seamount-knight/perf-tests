package main

import (
	"fmt"
	"net/http"
	"time"

	"k8s.io/perf-tests/clusterloader2/pkg/flags"
)

var (
	url = ""
)

func initFlags() {
	flags.StringVar(&url, "url", "", "")
}

func main() {
	flag := false

	initFlags()
	flags.Parse()

	fmt.Println("url: ", url)
	for !flag {
		if checkServer(url) {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		now := time.Now()
		checkServer(url)
		duration := time.Since(now)
		fmt.Println("time: ", duration)
		flag = true
	}

}

func checkServer(url string) bool {
	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
		}
		if resp.StatusCode == http.StatusOK {
			return true
		}
		fmt.Println(resp.StatusCode)
		time.Sleep(500 * time.Millisecond)
	}

}
