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
}
