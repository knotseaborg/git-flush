// Package flush is the backend of git-flush
package flush

import (
	utils "github.com/gitflush/utils"
)

var (
	logger = utils.InitLogger()
	config = utils.InitConfig()
)

func init() {
	err := config.Load()
	if err != nil {
		logger.Error("Failed to initialize flush")
	}
	if config.APIKey == "" {
		logger.Error("API Key not set in config file. Use `git-flush --config` to set your API Key")
	}
	// if err := godotenv.Load(); err != nil {
	// 	logger.Panic("No .env file found, or the .env file is of incorrect format")
}
