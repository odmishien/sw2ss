package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
	Sheet SheetConfig
}

type SheetConfig struct {
	SpreadSheetID string `toml:"spreadSheetID"`
	UpdateCell    string `toml:"updateCell"`
}

func readConfig() (Config, error) {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func getDuration(start time.Time, end time.Time) string {
	duration := end.Sub(start)
	hours := int(duration.Hours()) % 24
	mins := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d:%d:%d\n", hours, mins, secs)
}

func main() {
	start := time.Now()
	fmt.Printf("\x1b[36m%s\x1b[0m", "press Enter to stop your stopwatch!\n")
	bufio.NewScanner(os.Stdin).Scan()
	end := time.Now()
	duration := getDuration(start, end)
	fmt.Printf("result: %s \n", duration)

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	ctx := context.Background()
	srv, err := sheets.New(config.Client(ctx))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	conf, err := readConfig()
	if err != nil {
		log.Fatalf("Unable to read config.toml: %v", err)
	}
	spreadSheetID := conf.Sheet.SpreadSheetID
	updateCell := conf.Sheet.UpdateCell
	updateValue := &sheets.ValueRange{
		Values: [][]interface{}{
			{duration},
		},
	}
	_, err = srv.Spreadsheets.Values.Update(spreadSheetID, updateCell, updateValue).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}
	fmt.Printf("Success to update cell!\nIf you'd like to confirm the sheet, access:\nhttps://docs.google.com/spreadsheets/d/%s/edit#gid=0\n", spreadSheetID)
}
