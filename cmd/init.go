// Package cmd is the front-end of git-flush
package cmd

import (
	"fmt"
	"os"

	"github.com/gitflush/flush"
	"github.com/gitflush/utils"
	"github.com/spf13/cobra"
)

var (
	logger = utils.InitLogger()

	config        = utils.InitConfig()
	commitMessage string
	configMode    bool
	dryRun        bool
)

var rootCmd = &cobra.Command{
	Use:   "git-flush",
	Short: "Make commits with poop jokes!💩",
	Long: `git-flush is the equivalent of "git commit"
Commits are like pooping, so do it as frequently as you can for healthy code reviews and hilarious toilet humour!💩`,
	Run: func(cmd *cobra.Command, args []string) {
		// User wants to edit config
		if configMode {
			config.Edit()
			return
		}
		// User wants to test git-flush
		if dryRun {
			if err := joke(); err != nil {
				logger.Error("Dry run failed! Better check the plumbing🛠️", err)
			}
			return
		}
		// User wants to commit
		if commitMessage == "" {
			logger.Error("Oops... Commit message missing! Looks like you forgot to wipe the slate clean💩")
		} else {
			err := commitAndJoke(commitMessage)
			if err != nil {
				logger.Error("Flush failed! Give it another push🤗")
			}
		}
	},
}

func joke() error {
	repo := flush.InitWrapper()
	joker := flush.InitJoker()

	diff, err := repo.GetDiff()
	if err != nil {
		logger.Error("Looks like there's nothing to flush!😢")
		return err
	}
	joker.MakeJoke(diff)
	return nil
}

func commitAndJoke(message string) error {
	repo := flush.InitWrapper()
	joker := flush.InitJoker()

	diff, err := repo.GetDiff()
	if err != nil {
		logger.Error("Looks like there's nothing to flush!😢")
		return err
	}

	resp, err := repo.Commit(message)
	if err != nil {
		logger.Error("Looks like there's nothing to flush!😢")
		return err
	}
	fmt.Println(resp)
	joker.MakeJoke(diff)

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("An unknown clog occurred! Better check the plumbing🛠️", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&dryRun, "dryrun", "d", false, "dry run without commit")
	rootCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "commits with message")
	rootCmd.Flags().BoolVarP(&configMode, "config", "c", false, "edit config file")
}
