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

/* TrigerStat holds the information about about screens
 * triggered or not.
 * If triggered is true don't triggered. this screen is
 * already triggered.
 * If triggered is false */
type TriggerStat struct {
	Screens []ScreensTriggerStat `json:"Screens"`
}

type ScreensTriggerStat struct {
	Name      string `json:"name"`
	Triggered bool   `json"triggered"`
}

type config struct {
	Interval   string   `json:"interval"`
	Urls       []string `json:"urls"`
	TrigerFunc string   `json:"triger"`
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
var TriggerStructs []TriggerStat

func main() {
	flag.StringVar(&options.configFileName, "config", "", "name of config file json format")
	flag.Parse()

	/* --config argument is required.
	 * check --config has been set. */
	if len(os.Args) < 2 {
		log.Fatal("No args set. At least set --config <config-filename>")
	}

	readConfig(&cfg, options.configFileName)
	updateTriggerStats(&cfg)

	go runThread()

	quit := make(chan bool)
	<-quit
}

func runThread() {
	for {
		// read the config.
		readConfig(&cfg, options.configFileName)

		// check every ip that is given in the config file.
		checkScreens(&cfg)

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
func checkScreens(cfg *config) {
	var screenStats []parser.ScreenStat

	for _, value := range cfg.Urls {
		res, err := parser.GetScreenStatResponse(value)
		if err != nil {
			log.Printf("Error while getting screen stats: %s\n", err)
		}
		screenStats = append(screenStats, res)
	}

	// check screens are up or not
	for i, _ := range screenStats {
		log.Println("screen", screenStats)
		log.Println("triggerScreen", TriggerStructs)
		screen := screenStats[i].Screens
		triggerScreen := TriggerStructs[i].Screens

		for j, _ := range screen {
			log.Printf("Screen stats:\n\t\t%s is %t\n", screen[j].Name, screen[j].Up)

			// the screen went offline and message hasn't been send yet
			if !screen[j].Up && triggerScreen[j].Triggered == true {
				triggerFunc(cfg, screen[j].Name)
				triggerScreen[j].Triggered = false
			}

			// the screen went online(was offline) and we can send message
			if screen[j].Up && triggerScreen[j].Triggered == false {
				log.Printf("%s is online. Setting trigger to true.", screen[j].Name)
				triggerScreen[j].Triggered = true
			}
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

func updateTriggerStats(cfg *config) {
	TriggerStructs = []TriggerStat{}

	var triggerStructs []TriggerStat

	for _, value := range cfg.Urls {
		screenRes, err := parser.GetScreenStatResponse(value)
		if err != nil {
			log.Fatal(err)
		}

		var screenTriggerStats []ScreensTriggerStat
		for _, value := range screenRes.Screens {
			screensStruct := ScreensTriggerStat{
				Name:      value.Name,
				Triggered: true,
			}
			screenTriggerStats = append(screenTriggerStats, screensStruct)
		}

		triggerStruct := TriggerStat{
			Screens: screenTriggerStats,
		}

		triggerStructs = append(triggerStructs, triggerStruct)
	}
	TriggerStructs = triggerStructs
}
