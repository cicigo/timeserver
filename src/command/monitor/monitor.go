package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var targets []string

var result map[string]map[string][]Sample

type Sample struct {
	Time time.Time
	Value int
}

func monitor(targets []string, sampleIntervalSec int) {
	
	ticker := time.NewTicker(time.Duration(sampleIntervalSec) * time.Second)

	for {
		select {
		case <-ticker.C:
			for _, target := range targets {
				currentTime := time.Now()
				monitorResult := requestMonitor(target)
				if monitorResult == nil {
					break
				}
				for monitorName, monitorValue := range monitorResult {
					if _, ok := result[target]; !ok {
						result[target] = make(map[string][]Sample)
					}
					if _, ok := result[target][monitorName]; !ok {
						result[target][monitorName] = make([]Sample, 0)
					}
					sample := Sample{
						Time : currentTime,
						Value : monitorValue,
					}
					
					result[target][monitorName] = append(result[target][monitorName], sample)
				
				}
			}
			
		}
	}
}

func requestMonitor(target string) map[string]int {
	client := &http.Client{}
	url := fmt.Sprintf("%s/monitor", target)

	r, err := client.Get(url)

	if err != nil {
		fmt.Printf("Get request failed: %s", err)
		return nil
	}

	defer r.Body.Close()

	jsonBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("Read response failed: %s", err)
		return nil
	}

	counters := make(map[string]int)

	err = json.Unmarshal(jsonBytes, &counters)
	if err != nil {
		fmt.Printf("Unmarshalling json failed, %s", err)
		return nil
	} else {
		return counters
	}

}

func init() {
	result = make(map[string]map[string][]Sample)
}

func main() {

	targetsStr := flag.String("targets", "", "comma-separated list of URLs to monitor")
	sampleIntervalSec := flag.Int("sample-interval-sec", 10, "sample interval in seconds")
	runTimeSec := flag.Int("runtime-sec", 60, "monitor duration")

	flag.Parse()

	if *targetsStr == "" {
		fmt.Printf("targets is empty\n")
		return
	}
	var targets []string
	for _, target := range strings.Split(*targetsStr, ",") {
		target = strings.TrimSpace(target)
		if target != "" {
			targets = append(targets, target)
		}
	}
	if len(targets) == 0 {
		fmt.Printf("targets is empty\n")
		return
	}

	if *sampleIntervalSec <= 0 {
		fmt.Printf("sample interval sec must be greater than 0\n")
		return
	}

	if *runTimeSec <= 0 {
		fmt.Printf("runtime duration must be greater than 0\n")
		return
	}

	go monitor(targets, *sampleIntervalSec)
	time.Sleep(time.Duration(*runTimeSec) * time.Second)
	json, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		fmt.Printf("marshalling monitor result failed: %s\n", err)
	} else {
		fmt.Println(string(json[:]))
	}

}
