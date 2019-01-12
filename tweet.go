package main

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
)

type tweet struct {
	twitter.Tweet
}

func (t *tweet) HandleEntities(dir string) error {
	if t.ExtendedEntities != nil {
		for _, m := range t.ExtendedEntities.Media {
			switch m.Type {
			case "video":
				for _, v := range m.VideoInfo.Variants {
					if strings.Contains(v.URL, "mp4") {
						err := DownloadFile(v.URL, dir, "video.mp4")
						if err != nil {
							return err
						}
					}
				}
			case "photo":
				err := DownloadFile(m.MediaURL, dir, "photo.png")
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
