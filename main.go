package main

import (
	"AnimeQuote/providers"
	"fmt"
)

func main() {
	twitter := new(providers.Twitter)
	twitter = twitter.Init("", "", "", "", "https://api.twitter.com")


	media, err := twitter.UploadImage("https://static.wikia.nocookie.net/rezero/images/0/02/Echidna_Anime_PV_2.png/revision/latest?cb=20200611134057")
	if err != nil {
		panic(err)
	}

	fmt.Println(media)
}