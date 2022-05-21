package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gerifield/channel-checker/token"
)

type StreamInfo struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	TagIds       []string  `json:"tag_ids"`
	IsMature     bool      `json:"is_mature"`
}

func main() {
	channelsName := flag.String("channels", "gerifield,gibbonrike,marinemammalrescue", "Twitch channels name to check")

	clientID := flag.String("clientID", "", "Twitch App ClientID")
	clientSecret := flag.String("clientSecret", "", "Twitch App clientSecret")
	flag.Parse()

	tl := token.New(*clientID, *clientSecret)
	log.Println("Fetching token")
	token, err := tl.Get()
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.twitch.tv/helix/streams", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Client-Id", *clientID)

	params := req.URL.Query()
	for _, v := range strings.Split(*channelsName, ",") {
		params.Add("user_login", v)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var r struct {
		StreamInfos []StreamInfo `json:"data"`
		// Pagination too
	}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Println(err)
		return
	}

	for _, d := range r.StreamInfos {
		fmt.Printf("%s (%d) %s\n", d.UserName, d.ViewerCount, d.Type)
	}

}
