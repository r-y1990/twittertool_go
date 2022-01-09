package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
	yaml "gopkg.in/yaml.v2"
)

type TwitterSetting struct {
	ConsumerApiKey    string `yaml:"CONSUMER_API_KEY"`
	ConsumerSecretKey string `yaml:"CONSUMER_SECRET_KEY"`
	AccessApiKey      string `yaml:"ACCESS_API_KEY"`
	AccessSecretKey   string `yaml:"ACCESS_SECRET_KEY"`
}

func main() {
	//　APIツールの設定
	twitter := TwitterSetting{}
	conf, _ := os.ReadFile("config.yaml")
	yaml.Unmarshal(conf, &twitter)
	anaconda.SetConsumerKey(twitter.ConsumerApiKey)
	anaconda.SetConsumerSecret(twitter.ConsumerSecretKey)

	client := anaconda.NewTwitterApi(twitter.AccessApiKey, twitter.AccessSecretKey)
	v := url.Values{}
	v.Set("count", "50")
	next_cursor := "-1"

	file, err := os.OpenFile("data.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 全フォロワーサーチ
	for {
		v.Set("cursor", next_cursor)
		c, err := client.GetFollowersList(v)

		fmt.Println(err)

		for _, usr := range c.Users {
			writer.Write([]string{usr.IdStr, usr.Name, usr.ScreenName})
		}

		next_cursor = c.Next_cursor_str
		if err != nil || next_cursor == "0" {
			break
		}
	}
}
