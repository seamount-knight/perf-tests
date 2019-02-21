package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
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

	// go func() {
	// 	in.WriteString("kubectl  get node 10.21.128.13 | grep 10.21.128.13 | awk {'print $2'}\n")
	// }()
	flag1 := false
	flag2 := false

	var now time.Time
	var dura time.Duration

	for !flag1 {
		in := bytes.NewBuffer(nil)
		cmd := exec.Command("sh")
		var out bytes.Buffer
		cmd.Stdout = &out //输出

		cmd.Stdin = in
		in.WriteString("kubectl  get node 10.21.128.13 | grep 10.21.128.13 | awk {'print $2'}\n")

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		status := strings.TrimSpace(out.String())
		fmt.Printf("node status: %q\n", status)
		if status == "Ready" && !flag1 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		fmt.Println(status)
		if status == "NotReady" && !flag2 {
			flag2 = true
			now = time.Now()
		}
		if status == "Ready" && !flag2 {
			time.Sleep(500 * time.Millisecond)
			dura = time.Since(now)
			fmt.Println("dura", dura)
			flag1 = true
			return
		}
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
			fmt.Println("ok")
			return true
		}
		fmt.Println(resp.StatusCode)
		time.Sleep(500 * time.Millisecond)
	}

}
