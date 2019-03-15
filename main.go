package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"./parser"
)

var options struct {
	configFileName string
}

type config struct {
	Interval   string   `json:"interval"`
	ServerIPs  []string `json:"serverIPs"`
	TrigerFunc string   `json:"trigerFunc"`
	TrigerType string   `json:"trigerType"`
}

func readConfig(cfg *config, configFileName string) {
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("File error: ", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		log.Fatal("Config error: ", err.Error())
	}
}

var cfg config

func main() {
	flag.StringVar(&options.configFileName, "config", "", "name of config file json format")
	flag.Parse()

	/* --config argument is required.
	 * check --config has been set. */
	if len(os.Args) < 2 {
		log.Fatal("No args set. At least set --config <config-filename>")
	}
	go runThread()

	quit := make(chan bool)
	<-quit
}

func runThread() {
	for {
		/* read the config. */
		readConfig(&cfg, options.configFileName)

		/* check every ip that is given in the config file. */
		for _, value := range cfg.ServerIPs {
			url := "http://" + value + "/api/screens"
			checkScreens(&cfg, url)
		}
		interval, err := time.ParseDuration(cfg.Interval)
		if err != nil {
			log.Printf("Error while parsing interval from config: %s\n", err)
		}
		time.Sleep(interval)
	}
}

/* checkScreens gets screen stats while calling GetScreenStatResponse
 * with given url.
 * Checks each of the screen is up or not.
 * If screen is not up call triggerFunc. */
func checkScreens(cfg *config, url string) {
	log.Printf("Getting response from: %s\n", url)

	res, err := parser.GetScreenStatResponse(url)
	if err != nil {
		log.Printf("Error when getting screen stat: %s\n", err)
	}

	/* check screen is up or not */
	for _, value := range res.Screens {
		if !value.Up {
			triggerFunc(cfg, value.Name)
		}
	}
}

/* triggerFunc runs a bash script
 * that is given in the config file  from trigerFunc key*/
func triggerFunc(cfg *config, screenName string) {
	cmd := exec.Command(cfg.TrigerType, cfg.TrigerFunc)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error when executing trigger func: %s\n", err)
	}

	log.Printf("%s is not running. Has been send trigger func.\n", screenName)
}
