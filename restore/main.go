package main

import (
	"fmt"
	"net/http"
	"time"

	"k8s.io/perf-tests/clusterloader2/pkg/flags"
)

var (
	url string
)

func initFlags() {
	flags.StringVar(&url, "url", "", "")
}

func main() {
	initFlags()
	now := time.Now()
	checkServer(url)
	duration := time.Since(now)
	fmt.Println("time: ", duration)
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
