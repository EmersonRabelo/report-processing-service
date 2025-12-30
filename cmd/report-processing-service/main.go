package main

import (
	"fmt"

	"github.com/EmersonRabelo/report-processing-service/internal/config"
)

var setting config.SettingProvider

func init() {
	fmt.Println("Application initializing...")

	setting = config.GetSetting()

	config.InitDatabase()

	fmt.Println("Initialized.")
}

func main() {
	fmt.Println("Hello Wordl!")
}
