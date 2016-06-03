package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	strID := os.Getenv("ChannelID")
	numID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		log.Fatal("Wrong environment setting about ChannelID")
	}

	bot, err = linebot.NewClient(numID, os.Getenv("ChannelSecret"), os.Getenv("MID"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	received, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, result := range received.Results {
		content := result.Content()
		if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
			text, err := content.TextContent()

			log.Println("INPUT = " + text.Text)

			var outputString = stackoverflow(text.Text)

			log.Println("OUTPUT = " + outputString)

			_, err = bot.SendText([]string{content.From}, outputString)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

//Items:
type jsonobject struct {
	Items []Item
}

//Item
type Item struct {
	Link  string `json:"link"`
	Title string `json:"title"`
}

func stackoverflow(input string) string {

	root := "http://api.stackexchange.com/2.2/similar"
	para := "?page=1&pagesize=1&order=desc&sort=relevance&site=stackoverflow&title=" + url.QueryEscape(input)

	stackoverflowEndPoint := root + para

	resp, err := http.Get(stackoverflowEndPoint)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var i jsonobject
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Println(err)
	}

	var ret string

	if len(i.Items) == 0 {
		ret = "No Data Found"
	} else {
		ret = html.UnescapeString(i.Items[0].Title) + " " + i.Items[0].Link
	}

	if len(ret) == 0 {
		ret = "No Data"
	}

	return ret
}
