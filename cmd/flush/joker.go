package flush

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	SystemPrompt = "You are an english programmer with highly suggestive toilet humour." +
		"Use the Queen's english to roast the diff and make a single wildest poop joke" +
		"according to the git diffs provided to you. Use less than 20 words"
)

type Joker interface {
	MakeJoke(string)
}

func InitJoker() Joker {
	client := LLMClient{}
	client.Init()

	// diffLimit is set to prevent a massive surge of tokens to the LLM
	tmp := os.Getenv("DIFF_LIMIT")
	diffLimit, err := strconv.Atoi(tmp)
	if err != nil {
		logger.Warn("Couldn't read `DIFF_LIMIT` from .env file, defaulting to 5000 characters")
		diffLimit = 5000
	}

	return &ToiletJoker{client, diffLimit}
}

type ToiletJoker struct {
	client    LLMClient
	diffLimit int
}

func (j *ToiletJoker) MakeJoke(gitDiff string) {
	if len(gitDiff) > j.diffLimit {
		fmt.Println("Damn.. That was a stink bomb! 😏 You just flushed a huge commit! 👏👏")
		return
	}

	startTime := time.Now()

	prompt := fmt.Sprintf(SystemPrompt+"\n%s", gitDiff)
	text, tokensUsed, err := j.client.Complete(prompt)
	if err != nil {
		logger.Error("Failed to make a joke")
		fmt.Println("Something went wrong! I'm constipated 😣")
	}

	endTime := time.Now()

	fmt.Println(text)
	fmt.Println("---")
	fmt.Println("Pooped", tokensUsed, "tokens 💩\nConspitated for", endTime.Sub(startTime).Seconds(), "seconds 😏")
	fmt.Println("---")
}
