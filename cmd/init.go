// Package cmd is the front-end of git-flush
package cmd

import (
	"fmt"
	"os"

	"github.com/gitflush/cmd/flush"
	"github.com/spf13/cobra"
)

var (
	logger = flush.InitLogger()
	repo   = flush.InitWrapper()
	joker  = flush.InitJoker()

	commitMessage string
)

var rootCmd = &cobra.Command{
	Use:   "git-flush",
	Short: "Make commits with poop jokes!💩",
	Long: `Commits are like pooping, do it as frequently as you can for regular reviews and hilarious toilet humour!💩
git-flush is the equivalent of "git commit -m". 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if commitMessage == "" {
			logger.Warn("Hey there, Pooper! I couldn't find a commit message🖕")
		} else {
			err := commitAndJoke(commitMessage)
			if err != nil {
				logger.Error("If you can't poop on the first try, try again!🤗")
			}
		}
	},
}

func commitAndJoke(message string) error {
	diff, err := repo.GetDiff()
	if err != nil {
		logger.Error("Looks like there's no poop to flush!😢")
		return err
	}

	resp, err := repo.Commit(message)
	if err != nil {
		logger.Error("Looks like there's no poop to flush!😢")
		return err
	}
	fmt.Println(resp)
	joker.MakeJoke(diff)

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("An unknown error occurred! Plumbing needs to be fixed! ", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "Commit message")
}
