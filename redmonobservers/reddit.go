package redmonobservers

import (
	"context"
	"fmt"
	"os"

	"github.com/reactivex/rxgo/v2"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type reminder struct {
	bot     reddit.Bot
	channel chan<- rxgo.Item
}

func (r *reminder) Post(p *reddit.Post) error {
	r.channel <- rxgo.Of(*p)

	return nil
}

func RedditProducer(filePath string) rxgo.Observable {
	f := func(_ context.Context, ch chan<- rxgo.Item) {

		fmt.Println("running the observable")
		if bot, err := reddit.NewBotFromAgentFile(filePath, 0); err != nil {
			fmt.Println("Failed to create bot handle", err)
		} else {
			cfg := graw.Config{
				Subreddits: []string{"mechmarket"},
			}

			handler := &reminder{
				bot:     bot,
				channel: ch,
			}

			if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
				fmt.Println("failed to start graw run", err)
				os.Exit(1)
			} else {
				fmt.Println("graw run failed", wait())
				os.Exit(1)
			}

		}

	}

	return rxgo.Defer([]rxgo.Producer{f})

}
