package main

import (
	"AnimeQuote/providers"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joao-1/animechan-go"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting bot...")

	err := godotenv.Load()
	if err != nil { panic(err) }
	
	twitter := new(providers.Twitter)
	twitter = twitter.Init(os.Getenv("CONSUMER_ID"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"), "https://api.twitter.com")

	client := http.Client{}

	animechan := animechan.Animechan{Client: &client}

	again := make(chan bool)

	go func() {
		for range again {
			fmt.Println("Trying again...")
			CallMakePost(twitter, client, animechan, again)
		}
	}()

	fmt.Println("Making first post...")
	CallMakePost(twitter, client, animechan, make(chan bool))

	for range time.Tick(time.Hour * 1) {
		fmt.Println("Making new post...")
		CallMakePost(twitter, client, animechan, again)
	}
}

func CallMakePost(twitter *providers.Twitter, client http.Client, animechan animechan.Animechan, again chan bool) {
	tweet, err := MakePost(twitter, client, animechan)
	if err != nil { again <- true }
	fmt.Println(tweet)
}

func MakePost(twitter *providers.Twitter, client http.Client, animechan animechan.Animechan) (providers.Tweet, error) {
	quote, err := animechan.Random().Only()
	if err != nil { return providers.Tweet{}, err }

	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	formatedQuote := fmt.Sprintf(`
	%q - %s, %s

	#anime #quotes #%s #%s 
	`, quote.Quote, quote.Character, quote.Anime, regex.ReplaceAllString(quote.Character, ""), regex.ReplaceAllString(quote.Anime, ""))

	anilistResponse, err := providers.GetAnilistCharacterImageURL(quote.Character, "https://graphql.anilist.co", client)
	if err != nil { return providers.Tweet{}, err }

	media, err := twitter.UploadImage(anilistResponse.Data.Character.Image.Large)
	if err != nil { return providers.Tweet{}, err }

	tweet, err := twitter.Tweet(providers.TweetParams{Body: formatedQuote, Image: media.Id})
	if err != nil { return providers.Tweet{}, err }
	
	return tweet, nil
}