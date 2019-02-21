package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"k8s.io/perf-tests/clusterloader2/pkg/flags"
)

var (
	apiurl      = ""
	scheduleurl = ""
	controlurl  = ""
)

func initFlags() {
	flags.StringVar(&apiurl, "apiurl", "", "")
	flags.StringVar(&scheduleurl, "scheduleurl", "", "")
	flags.StringVar(&controlurl, "controlurl", "", "")
}

func main() {
	check_master_node()
}

func check_master_node() {

	in := bytes.NewBuffer(nil)
	cmd := exec.Command("sh")
	cmd.Stdin = in
	go func() {
		in.WriteString("kubectl  get node 10.21.128.13 | grep 10.21.128.13 | awk {'print $2'}")
		in.WriteString("exit\n")
	}()

	// cmdStr := "kubectl  get node 10.21.128.13 | grep 10.21.128.13 | awk {'print $2'}"
	// cmd := exec.Command(cmdStr)

	var out bytes.Buffer
	cmd.Stdout = &out //输出

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}

func checkServer(url string) bool {
	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("ok")
			return true
		}
		fmt.Println(resp.StatusCode)
		time.Sleep(500 * time.Millisecond)
	}

}
