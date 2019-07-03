package wasgeit

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type Config struct {
	DropDb      bool
	SetupDb     bool
	LogLevel    string
	ChromiumUrl string
}

func GetConfiguration() Config {
	config := Config{}
	flag.BoolVar(&config.DropDb, "drop-db", false, "Whether to drop DB")
	flag.BoolVar(&config.SetupDb, "setup-db", false, "Whether to create DB tables")
	flag.StringVar(&config.LogLevel, "log-level", "Info", "Set log level")
	flag.StringVar(&config.ChromiumUrl, "chromium-host", "http://chromium:9222",
		"Host of chromium instance to connect to. Do not specify a path.")
	flag.Parse()
	return config
}

func ConfigureLogging(logLevel string) {
	log.SetFormatter(&log.TextFormatter{CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		tokens := strings.Split(frame.Function, ".")
		return strings.Join(tokens[len(tokens)-2:], "."), ""
	}})

	log.SetReportCaller(true)

	if level, err := log.ParseLevel(logLevel); err != nil {
		panic(err)
	} else {
		log.Info("Set log level to ", level)
		log.SetLevel(level)
	}
}
