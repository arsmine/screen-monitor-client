package parser

import (
	"encoding/json"
	"net/http"
)

/* OsStat holds the parsed json response from /api/osstats */
type OsStat struct {
	Timestamp  int    `json:"timestamp"`
	Hostname   string `json:"hostname"`
	MemoryStat struct {
		Total     int `json:"total"`
		Free      int `json:"free"`
		Available int `json:"available"`
		SwapTotal int `json:"swapTotal"`
		SwapFree  int `json:"swapFree"`
	} `json:"memoryStat"`
	CPUStat struct {
		User             int     `json:"user"`
		System           int     `json:"system"`
		Idle             int     `json:"idle"`
		UserPercentage   float64 `json:"userPercentage"`
		SystemPercentage int     `json:"systemPercentage"`
		IdlePercentage   float64 `json:"idlePercentage"`
	} `json:"cpuStat"`
	Uptime   int64 `json:"uptime"`
	DiskStat []struct {
		Name            string `json:"Name"`
		ReadsCompleted  int    `json:"ReadsCompleted"`
		WritesCompleted int    `json:"WritesCompleted"`
	} `json:"diskStat"`
	NetworkStat []struct {
		Name    string `json:"Name"`
		RxBytes int64  `json:"RxBytes"`
		TxBytes int64  `json:"TxBytes"`
	} `json:"networkStat"`
	LoadAverage struct {
		Avg1  float64 `json:"avg1"`
		Avg5  float64 `json:"avg5"`
		Avg15 float64 `json:"avg15"`
	} `json:"loadAverage"`
}

/* ScreenStat holds the parsed json response from /api/screens */
type ScreenStat struct {
	Timestamp uint64 `json:"timestamp"`
	Hostname  string `json:"hostname"`
	Screens   []struct {
		Name string `json:"name"`
		Up   bool   `json:"up"`
	} `json:"screens"`
}

/* GetOsStatResponse collects given url's api response,
 * parses to the ScreenStat struct and returns it. */
func GetOsStatResponse(url string) (*OsStat, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var stat OsStat

	err = json.NewDecoder(res.Body).Decode(&stat)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

/* GetScreenStatResponse collects given url's api response,
 * parses to the ScreenStat struct and returns it. */
func GetScreenStatResponse(url string) (*ScreenStat, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var stat ScreenStat
	err = json.NewDecoder(res.Body).Decode(&stat)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}
