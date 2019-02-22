package main

import (
	"bytes"
	"fmt"
	"log"
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
	// cmdStr := "curl -v https://10.21.128.13:6443/healthz --insecure | grep ok\n"
	// checkMasterSvc(cmdStr)
	check_node_ready()
}

func check_node_ready() {

	var now time.Time
	var dura time.Duration
	now = time.Now()

	for {
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

		if status == "Ready" {
			dura = time.Since(now)
			return
		}

		if status == "NotReady" {
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("duration: ", dura)
}

func checkMasterSvc(cmdstr string) {
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
		in.WriteString(cmdstr)

		err := cmd.Run()
		if err != nil && !flag2 {
			fmt.Println("err: ", err)
			time.Sleep(500 * time.Millisecond)
			flag2 = true
			now = time.Now()
			// continue
		}

		status := strings.TrimSpace(out.String())
		fmt.Printf("node status: %q\n", status)
		if status == "ok" && !flag1 && !flag2 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		// fmt.Println(status)
		// if status == "" && !flag2 {
		// 	flag2 = true
		// 	now = time.Now()
		// }
		if status == "ok" && flag2 {
			time.Sleep(500 * time.Millisecond)
			dura = time.Since(now)
			fmt.Println("dura", dura)
			flag1 = true
			return
		}
	}
}

// func checkMasterServer(host string) bool {
// 	for {
// 		resp, err := http.Get(url)
// 		if err != nil {
// 			fmt.Println(err)
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 		if resp.StatusCode == http.StatusOK {
// 			fmt.Println("ok")
// 			return true
// 		}
// 		fmt.Println(resp.StatusCode)
// 		time.Sleep(500 * time.Millisecond)
// 	}

// }
