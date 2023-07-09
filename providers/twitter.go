package providers

import (
	"AnimeQuote/helpers"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

type authorizer struct{}
func (a *authorizer) Add(req *http.Request) {}

type TweetParams struct {
	Body string
	Image string
}

type Media struct {
	Id string
}

type CreateMediaResponse struct {
	MediaID string `json:"media_id_string"`
}

type Tweet struct {
	Data struct {
		Id string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
	ImageUrl string
}

type TwitterProvider interface {
	Tweet(TweetParams) (Tweet, error)
	UploadImage(image string) (Media, error)
}

type Twitter struct {
	client *twitter.Client
}

func (t *Twitter) Init(consumerId, consumerSecret, accessToken, accessTokenSecret, apiHost string) (*Twitter) {
	config := oauth1.NewConfig(consumerId, consumerSecret)
	httpClient := config.Client(oauth1.NoContext, &oauth1.Token{Token: accessToken, TokenSecret: accessTokenSecret})
	
	t.client = &twitter.Client{
		Authorizer: &authorizer{},
		Client:     httpClient,
		Host:       apiHost,
	}

	return t
}

func (t *Twitter) Tweet(params TweetParams) (Tweet, error) {
	req := twitter.CreateTweetRequest{
		Text: params.Body,
	}

	if params.Image != "" {
		req.Media = &twitter.CreateTweetMedia{IDs: []string{params.Image}}
	}

	tweetResponse, err := t.client.CreateTweet(context.Background(), req)
	if err != nil { return Tweet{}, err} 

	enc, err := json.MarshalIndent(tweetResponse, "", "    ")
	if err != nil { return Tweet{}, err}

	fmt.Println(string(enc))
	var tweet Tweet
	errParse := json.Unmarshal(enc, &tweet)
	if errParse != nil { return Tweet{}, errParse }
	
	tweet.ImageUrl = helpers.ExtractTwitterImageURL(tweet.Data.Text)

	fmt.Println(tweet.ImageUrl)
	return tweet, nil
}

func (t *Twitter) UploadImage(image string) (Media, error) {
	imageBase64, err := helpers.DownloadImage(image)
	if err != nil { return Media{}, err }

	data := url.Values{}
	data.Set("media_data", imageBase64)

	req, err := http.NewRequest("POST", regexp.MustCompile("api").ReplaceAllString(t.client.Host, "upload") + "/1.1/media/upload.json", strings.NewReader(data.Encode()))
	if err != nil { return Media{}, err }
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.client.Client.Do(req)
	if err != nil { return Media{}, err }

	defer res.Body.Close()

	var CreateMediaResponse CreateMediaResponse
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil { return Media{}, err }	
	errParse := json.Unmarshal(bodyBytes, &CreateMediaResponse)
	if errParse != nil { return Media{}, errParse }	

	return Media{Id: CreateMediaResponse.MediaID}, nil
}