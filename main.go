// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
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
	stackoverflowEndPoint := "http://api.stackexchange.com/2.2/search?order=desc&sort=activity&site=stackoverflow&intitle=" + url.QueryEscape(input)

	log.Println(stackoverflowEndPoint)

	resp, err := http.Get(stackoverflowEndPoint)
	if err != nil {
		log.Println("RESPONSE ERROR")
		log.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("BODY ERROR")
		log.Println(err)
	}

	var i jsonobject
	err = json.Unmarshal(body, &i)
	if err != nil {
		log.Println("JSON ERROR")
		log.Println(err)
	}

	var ret = i.Items[0].Title + " " + i.Items[0].Link
	if len(ret) == 0 {
		ret = "No Data"
	}

	return ret
}
