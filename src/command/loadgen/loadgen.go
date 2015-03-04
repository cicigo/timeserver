package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"utils"
)

type LoadGenConfig struct {
	rate      int
	burst     int
	timeoutMs int
	runtime   int
	url       string
}

var totalCounter = utils.NewCounter("Total")
var counter100 = utils.NewCounter("100s")
var counter200 = utils.NewCounter("200s")
var counter300 = utils.NewCounter("300s")
var counter400 = utils.NewCounter("400s")
var counter500 = utils.NewCounter("500s")
var errorsCounter = utils.NewCounter("Errors")

var config LoadGenConfig

func incStatusCodeCounter(statusCode int) {
	switch statusCode {
	case 100:
		counter100.Incr(1)
	case 200:
		counter200.Incr(1)
	case 300:
		counter300.Incr(1)
	case 400:
		counter400.Incr(1)
	case 500:
		counter500.Incr(1)
	default:
		errorsCounter.Incr(1)

	}
}

func sendRequest(config LoadGenConfig) {
	timeout := make(chan bool)
	go func() {
		time.Sleep(time.Duration(config.timeoutMs) * time.Millisecond)
		timeout <- true
	}()

	statusCode := make(chan int)
	error := make(chan error)
	go func() {
		client := &http.Client{}
		response, err := client.Get(config.url)
		if err != nil {
			error <- err
		} else {
			defer response.Body.Close()
			statusCode <- response.StatusCode / 100 * 100
		}
	}()

	select {
	case <-timeout:
	case <-error:
		errorsCounter.Incr(1)
	case code := <-statusCode:
		incStatusCodeCounter(code)
	}
}

func genLoad(config LoadGenConfig) {
	tickDuration := config.burst * 1000000 / config.rate

	ticker := time.NewTicker(time.Duration(tickDuration) * time.Microsecond)

	for {
		select {
		case <-ticker.C:
			for i := 0; i < config.burst; i++ {
				totalCounter.Incr(1)
				go sendRequest(config)
			}
		}
	}
}

func main() {

	config = LoadGenConfig{}

	flag.IntVar(&config.rate, "rate", 1, "Rate")
	flag.IntVar(&config.burst, "burst", 1, "Burst")
	flag.IntVar(&config.runtime, "runtime", 60, "Run time")
	flag.IntVar(&config.timeoutMs, "timeout-ms", 1000, "time out in millisecond")
	flag.StringVar(&config.url, "url", "", "url to load test")

	flag.Parse()

	if config.url == "" {
		fmt.Printf("Url is empty.\n")
		return
	} else {
		fmt.Printf("Load generating to hit url %s\n", config.url)
	}

	go genLoad(config)
	time.Sleep(time.Duration(config.runtime) * time.Second)
	counters := utils.DumpCounter()
	for name, value := range counters {
		fmt.Printf("%s: %v\n", name, value)
	}
}
