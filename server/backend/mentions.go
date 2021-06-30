package backend

import (
	"context"
	"encoding/json"
	_ "fmt"

	"github.com/fieldkit/cloud/server/data"
)

type Mention struct {
	UserID int32
	PostID int64
}

func findTipTapMentions(ctx context.Context, value interface{}) []int32 {
	ids := make([]int32, 0)
	switch vv := value.(type) {
	case map[string]interface{}:
		if nodeType, ok := vv["type"]; ok {
			if nodeType == "mention" {
				if attrs, ok := vv["attrs"].(map[string]interface{}); ok {
					if id, ok := attrs["id"].(float64); ok {
						return []int32{int32(id)}
					}
				}
			}
		}
		for _, child := range vv {
			ids = append(ids, findTipTapMentions(ctx, child)...)
		}
	case []interface{}:
		for _, child := range vv {
			ids = append(ids, findTipTapMentions(ctx, child)...)
		}
	}

	return ids
}

func DiscoverMentions(ctx context.Context, postID int64, body string) ([]*Mention, error) {
	log := Logger(ctx).Sugar()

	var tiptap interface{}
	if err := json.Unmarshal([]byte(body), &tiptap); err != nil {
		log.Infow("non-json body, skipping mentions")
		return nil, nil
	}

	ids := findTipTapMentions(ctx, tiptap)

	log.Infow("mentions", "ids", ids)

	mentions := make([]*Mention, 0)

	for _, id := range ids {
		mentions = append(mentions, &Mention{
			UserID: id,
			PostID: postID,
		})
	}

	return mentions, nil
}

func NotifyMentions(mentions []*Mention) []*data.Notification {
	notifications := make([]*data.Notification, len(mentions))
	for i, mention := range mentions {
		notifications[i] = data.NewMentionNotification(mention.UserID, mention.PostID)
	}
	return notifications
}
