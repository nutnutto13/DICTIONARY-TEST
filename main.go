package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	channelSecret := "f8830b25d8c4f293448712ddf9c5cc12"
	channelAccessToken := "e0OFDDYPrDJZE4ZcUAJVAmIGxRy+ZhcYiSA9Yqe9Uk+M7dFRF4lKOtjDtTII/6r78hqlumtw/Z5zolTOf4589mXpdkTyqldbOntJrWOpe4SlVh44veIihUGGQvIpcfHV9D0HeXF7UTQfPIi1UqGw6QdB04t89/1O/w1cDnyilFU="
	log.Print("start......")
	bot, err := linebot.New(channelSecret, channelAccessToken)
	if err != nil {
		log.Fatal(err)
	}

	//webhook
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		//for loop suport word has multiple meanings
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				if message, ok := event.Message.(*linebot.TextMessage); ok {

					var newTextMessages []linebot.SendingMessage
					result, err := GetDictionary(message.Text)
					if err != nil {
						newTextMessages = append(newTextMessages, linebot.NewTextMessage("Sorry : "+err.Error()))
					} else {
						for i, item := range result {
							newTextMessages = append(newTextMessages, linebot.NewTextMessage("Meanings"+strconv.Itoa(i+1)+" : "+item.Meanings[0].Definitions[0].Definition+"\nSynonyms: "+strings.Join(item.Meanings[0].Definitions[0].Synonyms, " ")))
						}
					}

					_, err = bot.ReplyMessage(event.ReplyToken, newTextMessages...).Do()
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetDictionary(word string) ([]Result, error) {
	//validate input
	var validWord = regexp.MustCompile(`^[a-zA-Z]+$`)
	if !validWord.MatchString(word) {
		return []Result{}, errors.New("input invalid, input English only")
	}

	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word)
	resp, err := http.Get(url)
	if err != nil {
		return []Result{}, errors.New("http get error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Result{}, errors.New("http get error")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Result{}, errors.New("http get error")
	}

	//unmarshal dictionary output to model
	var result []Result
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return []Result{}, errors.New("unmarshal error")
	}

	return result, nil

}

// Model
type Definition struct {
	Definition string   `json:"definition"`
	Synonyms   []string `json:"synonyms"`
	Antonyms   []string `json:"antonyms"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
	Synonyms     []string     `json:"synonyms"`
	Antonyms     []string     `json:"antonyms"`
}

type Result struct {
	Word     string    `json:"word"`
	Meanings []Meaning `json:"meanings"`
}
