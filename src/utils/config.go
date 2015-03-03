package utils

import (
	"flag"
)

var config = new(Config)

func GetConfig() *Config {
	return config
}

type Config struct {
	Version            bool
	Port               uint
	AuthPort           uint
	AuthHost           string
	AuthTimeoutMs      float64
	AvgResponseMs      float64
	DeviationMs        float64
	CheckpointInterval int
	MaxInflight        int
	Templates          string
	Log                string
	DumpFile           string
}

func init() {
	flag.BoolVar(&config.Version, "v", false, "Display version number")
	flag.UintVar(&config.Port, "port", 8080, "Time server port number")
	flag.UintVar(&config.AuthPort, "authport", 7070, "Auth server port number")
	flag.StringVar(&config.AuthHost, "authhost", "", "Auth server host")
	flag.Float64Var(&config.AuthTimeoutMs, "authtimeout-ms", 1000.0, "Auth response timeout in milliseconds")
	flag.Float64Var(&config.AvgResponseMs, "avg-response-ms", 0.0, "Avg response in milliseconds")
	flag.Float64Var(&config.DeviationMs, "deviation-ms", 0.0, "Response deviation in milliseconds")
	flag.IntVar(&config.CheckpointInterval, "checkpoint-interval", 60, "Checkpoint save interval in seconds")
	flag.IntVar(&config.MaxInflight, "max-inflight", 0, "Maximum number of in-flight time requests the server can handle")
	flag.StringVar(&config.Templates, "templates", "templates", "Templates folder")
	flag.StringVar(&config.Log, "log", "etc/log.xml", "Log configuration file path")
	flag.StringVar(&config.DumpFile, "dumpfile", "", "auth info dump file")

	flag.Parse()
}
