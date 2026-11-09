package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Sn0wo2/go-afdian-api"
	"github.com/parquet-go/bitpack/unsafecast"
)

func main() {
	client := afdian.NewClient(&afdian.Config{
		UserID:   os.Getenv("USER_ID"),
		APIToken: os.Getenv("API_TOKEN"),
	})

	sponsor, err := client.QuerySponsor(1, 100)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(sponsor.RawResponse.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(unsafecast.String(body))
}
