package blur

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	http "github.com/bogdanfinn/fhttp"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

const (
	blurCollectionURL           = "https://core-api.prod.blur.io/v1/collections/%s"
	blurCollectionOwnershipsURL = "https://core-api.prod.blur.io/v1/collections/%s/ownerships"
)

type Client struct {
	tlsClient.HttpClient
}

func NewClient() (*Client, error) {
	options := []tlsClient.HttpClientOption{
		tlsClient.WithTimeoutSeconds(30),
		tlsClient.WithClientProfile(profiles.Chrome_105),
		tlsClient.WithNotFollowRedirects(),
	}

	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create custom tls http client: %w", err)
	}

	return &Client{
		client,
	}, err
}

type CollectionOwnerships struct {
	Ownerships []struct {
		OwnerAddress string `json:"ownerAddress"`
		NumberOwned  int    `json:"numberOwned"`
	} `json:"ownerships"`
}

func (c *Client) GetCollectionOwnerships(collectionAddress string) (CollectionOwnerships, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(blurCollectionOwnershipsURL, collectionAddress), nil)
	if err != nil {
		return CollectionOwnerships{}, err
	}

	req.Header = http.Header{
		"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"},
	}

	resp, err := c.Do(req)
	if err != nil {
		return CollectionOwnerships{}, err
	}

	defer resp.Body.Close()

	log.Println(fmt.Sprintf("status code: %d", resp.StatusCode))

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return CollectionOwnerships{}, err
	}

	data := new(CollectionOwnerships)
	if err := json.Unmarshal(readBytes, data); err != nil {
		return CollectionOwnerships{}, err
	}

	return *data, nil
}

type CollectionName struct {
	Collection struct {
		Name string `json:"name"`
	} `json:"collection"`
}

func (c *Client) GetCollectionNameByAddress(collectionAddress string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(blurCollectionURL, collectionAddress), nil)
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"},
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("status code: %d", resp.StatusCode))
	}

	readBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := new(CollectionName)
	if err := json.Unmarshal(readBytes, data); err != nil {
		return "", err
	}

	return data.Collection.Name, nil
}
