package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	http "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type blurCollectionHolders struct {
	Success    bool `json:"success"`
	Ownerships []struct {
		OwnerAddress string `json:"ownerAddress"`
		NumberOwned  int    `json:"numberOwned"`
	} `json:"ownerships"`
}

func getCollectionHolders(collectionAddress string) (*blurCollectionHolders, error) {
	jar := tlsClient.NewCookieJar()
	options := []tlsClient.HttpClientOption{
		tlsClient.WithTimeoutSeconds(30),
		tlsClient.WithClientProfile(profiles.Chrome_105),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://core-api.prod.blur.io/v1/collections/%s/ownerships", collectionAddress), nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"accept":          {"*/*"},
		"accept-language": {"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7"},
		"user-agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"},
		http.HeaderOrderKey: {
			"accept",
			"accept-language",
			"user-agent",
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Println(fmt.Sprintf("status code: %d", resp.StatusCode))

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := new(blurCollectionHolders)
	if err := json.Unmarshal(readBytes, data); err != nil {
		return nil, err
	}

	return data, nil
}
