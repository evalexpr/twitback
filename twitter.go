package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

// Client wraps a Twitter client
type Client struct {
	twitterClient *twitter.Client
	twitterUser   *twitter.User
}

// NewClient returns a Twitter client with oAuth embedded
func NewClient() Client {
	config := oauth1.NewConfig(cKey, cSecret)
	token := oauth1.NewToken(tToken, tSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return Client{twitterClient: twitter.NewClient(httpClient)}
}

// VerifyCredentials verifies that the given credentials were correct
func (c *Client) VerifyCredentials() error {
	params := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err := c.twitterClient.Accounts.VerifyCredentials(params)
	c.twitterUser = user

	return err
}

var (
	minID int64
)

// DownloadFavorites downloads the users favorites to the specified directory
func (c *Client) DownloadFavorites() error {
	params := &twitter.FavoriteListParams{
		UserID:          c.twitterUser.ID,
		Count:           200,
		IncludeEntities: twitter.Bool(true),
	}
	if minID > 0 {
		params.MaxID = minID
	}
	tweets, res, err := c.twitterClient.Favorites.List(params)
	if res.StatusCode != 200 || len(tweets) == 0 {
		return errors.New("could not fetch favourites")
	}
	if len(tweets) == 1 {
		minID = 0
		return nil
	}

	for _, t := range tweets {
		if minID == 0 {
			minID = t.ID
		}
		if t.ID < minID {
			minID = t.ID
		}

		format := "Mon Jan 2 15:04:05 +0000 2006"
		time, err := time.Parse(format, t.CreatedAt)
		if err != nil {
			return fmt.Errorf("could not parse time: %v", err)
		}
		dir := directory + "/" + time.Format("2006-01-02") + "-" + t.IDStr
		text := dir + "/tweet.txt"
		tw := tweet{t}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating dir for tweet %s failed: %v", dir, err)
			}
		} else {
			// Skip this tweet
			logrus.Debugf("Skipping tweet %d in folder %s", tw.ID, dir)
			continue
		}
		if _, err := os.Stat(text); os.IsNotExist(err) {
			f, err := os.Create(text)
			if err != nil {
				return fmt.Errorf("could not create tweet.txt: %v", err)
			}
			defer f.Close()

			if _, err := f.WriteString(t.Text); err != nil {
				return fmt.Errorf("could not write to tweet.txt: %v", err)
			}
		}
		tw.HandleEntities(dir)
	}

	c.DownloadFavorites()
	return err
}
