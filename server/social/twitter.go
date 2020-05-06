package social

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/fieldkit/cloud/server/backend"
	"github.com/fieldkit/cloud/server/social/twitter"
)

const (
	TwitterJSONSchema = `
{
	"id": "https://api.fieldkit.org/schemas/tweet.json",
	"type": "object",
	"properties": {
		"id_str": {
			"type": "string"
		},
		"text": {
			"type": "string"
		}
	}
}
`
)

type TwitterOptions struct {
	StreamOptions twitter.StreamOptions
	Backend       *backend.Backend
	Done          <-chan struct{}
}

func Twitter(options TwitterOptions) {
	stream := twitter.NewStream(options.StreamOptions)
	for {
		select {
		case <-options.Done:
			stream.Stop()
			return

		case tweet := <-stream.Tweets:
			tweetData := new(bytes.Buffer)
			if err := json.NewEncoder(tweetData).Encode(tweet); err != nil {
				log.Println("!", "tweet error:", err)
				continue
			}

			ctx := context.Background()
			_, err := options.Backend.ListTwitterAccountSourcesByAccountID(ctx, tweet.User.ID)
			if err != nil {
				log.Println("!", "twitter account sources error:", err)
				continue
			}
		}
	}
}
