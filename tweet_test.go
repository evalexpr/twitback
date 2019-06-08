package main

import (
	"testing"

	"github.com/dghubble/go-twitter/twitter"
)

func TestHandleEntities(t *testing.T) {
	type fields struct {
		Tweet twitter.Tweet
	}
	type args struct {
		dir string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		shouldErr bool
	}{
		{
			"No media returns nil",
			fields{twitter.Tweet{ExtendedEntities: &twitter.ExtendedEntity{Media: nil}}},
			args{"foo"},
			false,
		},
		{
			"Photo downloads the photo",
			fields{
				twitter.Tweet{
					ExtendedEntities: &twitter.ExtendedEntity{
						Media: []twitter.MediaEntity{
							{Type: "photo", MediaURL: "https://static.wilkins.tech/twitback/test.png"},
						},
					},
				},
			},
			args{"test_fixtures"},
			false,
		},
		{
			"Video downloads the video",
			fields{
				twitter.Tweet{
					ExtendedEntities: &twitter.ExtendedEntity{
						Media: []twitter.MediaEntity{
							{
								Type: "video",
								VideoInfo: twitter.VideoInfo{
									Variants: []twitter.VideoVariant{{URL: "https://static.wilkins.tech/twitback/test.mp4"}},
								},
							},
						},
					},
				},
			},
			args{"test_fixtures"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &tweet{
				Tweet: tt.fields.Tweet,
			}
			if err := tw.HandleEntities(tt.args.dir); (err != nil) != tt.shouldErr {
				t.Errorf("tweet.HandleEntities() failed, error: %v, shouldErr: %v", err, tt.shouldErr)
			}
		})
	}
}
