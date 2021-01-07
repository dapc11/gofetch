package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func getHostname() string {
	name, err := os.Hostname()
	check(err)

	return name
}

func getDistro() string {
	distro, err := exec.Command("lsb_release", "-a").Output()
	check(err)

	re := regexp.MustCompile("Description:(.*)")
	match := re.FindStringSubmatch(string(distro))
	return strings.TrimSpace(match[1])
}

func getUptime() string {
	uptime, err := ioutil.ReadFile("/proc/uptime")
	check(err)
	upSplit := strings.Split(string(uptime), ".")[0]
	fmt.Println("#####")

	now := time.Now()
	upInSeconds, err := strconv.Atoi(upSplit)
	check(err)

	after := now.Add(-time.Duration(upInSeconds) * time.Second)

	return formatSince(after)
}

func formatSince(t time.Time) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	ts := time.Since(t)
	sign := time.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	s := ts / time.Second
	ts = ts % time.Second
	return fmt.Sprintf("%dd%dh%dm%ds", d, h, m, s)
}

func getKernel() string {
	kernel, err := exec.Command("uname", "-mrs").Output()
	check(err)

	return strings.TrimSpace(string(kernel))
}

func getProductName() string {
	productName, err := ioutil.ReadFile("/sys/devices/virtual/dmi/id/product_name")
	check(err)

	return strings.TrimSpace(string(productName))
}

func main() {
	hostname := getHostname()
	distro := getDistro()
	kernel := getKernel()
	uptime := getUptime()

	fmt.Println(hostname)
	fmt.Println(distro)
	fmt.Println(uptime)
	fmt.Println(kernel)
	fmt.Println(getProductName())
}
