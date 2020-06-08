package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type Config struct {
	Sheet SheetConfig
}

type SheetConfig struct {
	SpreadSheetID string `toml:"spreadSheetID"`
	UpdateCell    string `toml:"updateCell"`
}

var home = os.Getenv("HOME")

const (
	configDir          = ".config"
	projectName        = "sw2ss"
	configFileName     = "config.toml"
	credentialFileName = "credentials.json"
)

func loadConfig() (Config, error) {
	var config Config
	confFilePath := filepath.Join(home, configDir, projectName, configFileName)
	_, err := toml.DecodeFile(confFilePath, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func loadCredential() ([]byte, error) {
	credentialFilePath := filepath.Join(home, configDir, projectName, credentialFileName)
	b, err := ioutil.ReadFile(credentialFilePath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getDuration(start time.Time, end time.Time) string {
	duration := end.Sub(start)
	hours := int(duration.Hours()) % 24
	mins := int(duration.Minutes()) % 60
	secs := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d:%d:%d\n", hours, mins, secs)
}

func getSheetClient(config *jwt.Config) (*sheets.Service, error) {
	ctx := context.Background()
	srv, err := sheets.New(config.Client(ctx))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv, nil
}

func getCommand(args []string) (string, error) {
	if len(args) == 0 {
		// もしコマンドが指定されない場合は start ってことにする
		return "start", nil
	}
	var commands = []string{"start", "config"}
	inputCommand := args[0]
	for _, v := range commands {
		if v == inputCommand {
			return inputCommand, nil
		}
	}
	return "", fmt.Errorf("unknown command: %s", inputCommand)
}

func isValidConfigItems(item string) bool {
	var configItem = []string{"spreadSheetID", "updateCell"}
	for _, v := range configItem {
		if v == item {
			return true
		}
	}
	return false
}

func start() {
	start := time.Now()
	fmt.Printf("\x1b[36m%s\x1b[0m", "press Enter to stop your stopwatch!\n")
	bufio.NewScanner(os.Stdin).Scan()
	end := time.Now()
	duration := getDuration(start, end)
	fmt.Printf("result: %s \n", duration)

	credential, err := loadCredential()
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.JWTConfigFromJSON(credential, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	srv, err := getSheetClient(config)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	conf, err := loadConfig()
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

func main() {
	flag.Parse()
	command, err := getCommand(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	switch command {
	case "start":
		start()
	case "config":
		// TODO: configというコマンドでtomlファイルをいじれるようにする？
		return
	}

}
