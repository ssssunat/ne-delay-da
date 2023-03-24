package main

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"bytes"
	"strconv"
)


type Update struct {
	UpdateId int `json:"update_id"`	
	Message Message `json:"message"`
}

type Message struct {
	Chat Chat `json:"chat"`
	Text string	`json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId	int	`json:"chat_id"`
	Text	string	`json:"text"`
}
//точка входа программы
func main() {
	botToken := "5695972547:AAEHier6gjQqr1IWOZL31vx3du47Aw__lXc"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken

	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("smth wrong ")
		}
		
		for _, update := range updates {
			err = respond(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}


//запрос обновлений
func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

//ответ на обновления
func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text + " , не делай даа"

	buf, err := json.Marshal(botMessage)

	if err != nil {
		return err
	}
	_, err = http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))
	if  err != nil {
		return err
	}
	return nil
}