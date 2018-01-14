package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	LOGGER = log.New(os.Stdout, fmt.Sprintf("[%s %s] ", APP_NAME, VERSION), log.Ltime)
	config_path := flag.String("config", "", "Path to configuration file")
	flag.Parse()
	if *config_path != "" {
		cfg, err := ParseConfig(*config_path)
		if err != nil {
			fmt.Println(fmt.Sprintf("[ERROR] %v", err))
		} else {
			LOG(cfg, fmt.Sprintf("Starting %s %s", APP_NAME, VERSION))
			StartServer(LOGGER, cfg)
		}
	} else {
		fmt.Println("[ERROR] Missing arguments\nUsage:")
		flag.PrintDefaults()
	}
}
