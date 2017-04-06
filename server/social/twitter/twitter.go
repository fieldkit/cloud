package twitter

import (
	"log"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
)

func Sleep(duration time.Duration, done <-chan struct{}) {
	timer := time.NewTimer(duration)
	select {
	case <-timer.C:
	case <-done:
		timer.Stop()
	}
}

// UserLister wraps the UserList method.
type UserLister interface {
	UserList() (ids []int64, err error)
}

type UserCredentialer interface {
	UserCredentials(id int64) (accessToken, accessSecret string, err error)
}

type StreamOptions struct {
	ConsumerKey      string
	ConsumerSecret   string
	Domain           string
	UserLister       UserLister
	UserCredentialer UserCredentialer
}

type Stream struct {
	Tweets chan *twitter.Tweet

	done              chan struct{}
	waitgroup         *sync.WaitGroup
	userLister        UserLister
	userCredentialler UserCredentialer
	oauth1Config      *oauth1.Config
}

func NewStream(options StreamOptions) *Stream {
	s := &Stream{
		Tweets:            make(chan *twitter.Tweet),
		userLister:        options.UserLister,
		userCredentialler: options.UserCredentialer,
		done:              make(chan struct{}),
		waitgroup:         &sync.WaitGroup{},
		oauth1Config: &oauth1.Config{
			ConsumerKey:    options.ConsumerKey,
			ConsumerSecret: options.ConsumerSecret,
			CallbackURL:    "https://api." + options.Domain + "/twitter/callback",
			Endpoint:       oauth1Twitter.AuthorizeEndpoint,
		},
	}

	s.waitgroup.Add(1)
	go s.listen()
	return s
}

func (l *Stream) Stop() {
	close(l.done)
	l.waitgroup.Wait()
	close(l.Tweets)
}

func (l *Stream) userStream(id int64) (*twitter.Stream, error) {
	accessToken, accessSecret, err := l.userCredentialler.UserCredentials(id)
	if err != nil {
		return nil, err
	}

	client := twitter.NewClient(l.oauth1Config.Client(oauth1.NoContext, oauth1.NewToken(accessToken, accessSecret)))
	return client.Streams.User(&twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
		With:          "user",
	})
}

func handleStream(stream *twitter.Stream, demux twitter.Demux, done <-chan struct{}) bool {
	for {
		select {
		case <-done:
			stream.Stop()
			return true

		case message, ok := <-stream.Messages:
			if !ok {
				Sleep(5*time.Second, done)
				return false
			}

			log.Println(message)
			demux.Handle(message)
		}
	}

	return true
}

func (l *Stream) listenUser(id int64, done <-chan struct{}) {
	defer l.waitgroup.Done()

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		l.Tweets <- tweet
	}

	for {
		stream, err := l.userStream(id)

		if err != nil {
			log.Println("!", "stream error:", err)
			Sleep(time.Minute, done)
			continue
		}

		if handleStream(stream, demux, done) {
			return
		}
	}
}

func (l *Stream) listen() {
	log.Println("listen startup")
	defer l.waitgroup.Done()

	userStreams := map[int64]chan struct{}{}
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-l.done:
			log.Println("listen shutdown")

			// Stop the ticker.
			ticker.Stop()

			// Stop all user listeners.
			for _, done := range userStreams {
				close(done)
			}

			return

		case <-ticker.C:
			log.Println("reloading user list")

			// Get the current list of Twitter users.
			userIDs, err := l.userLister.UserList()
			if err != nil {
				log.Println("!", "list error:", err)
				continue
			}

			// Build a map.
			userIDsMap := map[int64]struct{}{}
			for _, userID := range userIDs {
				userIDsMap[userID] = struct{}{}
			}

			// Stop missing user listeners.
			for userID, done := range userStreams {
				if _, ok := userIDsMap[userID]; !ok {
					log.Println("-", userID)

					close(done)
					delete(userStreams, userID)
				}
			}

			for _, userID := range userIDs {

				// Check if a listener is already running for this user ID.
				if _, ok := userStreams[userID]; ok {
					continue
				}

				log.Println("+", userID)

				done := make(chan struct{})
				userStreams[userID] = done

				l.waitgroup.Add(1)
				go l.listenUser(userID, done)
			}
		}
	}
}

func Coordinates(tweet *twitter.Tweet) [2]float64 {
	if tweet.Coordinates != nil {
		return tweet.Coordinates.Coordinates
	}

	coordinates := [2]float64{}
	if tweet.Place == nil {
		return coordinates
	}

	if tweet.Place.BoundingBox == nil {
		return coordinates
	}

	if len(tweet.Place.BoundingBox.Coordinates) != 1 {
		return coordinates
	}

	points := tweet.Place.BoundingBox.Coordinates[0]
	for _, point := range points {
		coordinates[0] += point[0]
		coordinates[1] += point[1]
	}

	pointsCount := float64(len(tweet.Place.BoundingBox.Coordinates[0]))
	coordinates[0] /= pointsCount
	coordinates[1] /= pointsCount
	return coordinates
}
