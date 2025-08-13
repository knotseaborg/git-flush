// Package flush is the backend of git-flush
package flush

import "github.com/joho/godotenv"

var logger = InitLogger()

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Panic("No .env file found, or the .env file is of incorrect format")
	}
}
