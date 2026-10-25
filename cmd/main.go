package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Sn0wo2/go-afdian-api"
	"github.com/Sn0wo2/go-afdian-api/internal/helper"
)

func main() {
	config := &afdian.Config{
		UserID:   os.Getenv("USER_ID"),
		APIToken: os.Getenv("API_TOKEN"),

		/*WebHookListenAddr: ":9000",
		  WebHookPath:       "/WH",
		  WebHookCallback: func(p *payload.WebHook, errs ...error) {
		  	body, err := io.ReadAll(p.RawRequest.Body)
		  	if err != nil {
		  		panic(err)
		  	}
		  	fmt.Println(helper.BytesToString(body))
		  },*/
	}
	client := afdian.NewClient(config)

	sponsor, err := client.QuerySponsor(1, 100)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(sponsor.RawResponse.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(helper.BytesToString(body))

	/*if err := afdian.NewWebHook(client).Start(); err != nil {
		panic(err)
	}*/
}
