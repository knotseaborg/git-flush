package flush

import (
	"fmt"
	"time"
)

const (
	SystemPrompt = "You are a senior 'dad' programmer notorious for toilet humour. " +
		"Roast this git diff and make a single wild poop joke zinger " +
		"using less than 20 words"
)

type Joker interface {
	MakeJoke(string)
}

func InitJoker() Joker {
	client := InitLLMClient()

	// diffLimit is set to prevent a massive surge of tokens to the LLM
	return &ToiletJoker{client, config.DiffLimit}
}

type ToiletJoker struct {
	client    *LLMClient
	diffLimit int
}

func (j *ToiletJoker) MakeJoke(gitDiff string) {
	if len(gitDiff) > j.diffLimit {
		fmt.Println("Woah! That's one mega-dump!😲 You just flushed a huge commit! 👏👏")
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
