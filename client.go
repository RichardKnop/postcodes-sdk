package idealpostcodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Endpoint is the API base URI
const Endpoint = "https://api.ideal-postcodes.co.uk/v1"

// Client for https://ideal-postcodes.co.uk API
type Client struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

// NewClient returns new Client instance
func NewClient(endpoint string, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = new(http.Client)
	}
	return &Client{
		endpoint:   endpoint,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// GetPostcode returns list of addresses matching a postcode
func (c *Client) GetPostcode(postcode string) ([]*Address, error) {
	// Prepare a request
	payload := url.Values{}
	payload.Add("api_key", c.apiKey)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/postcodes/%s?%s", c.endpoint, postcode, payload.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the response data
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshall into our struct
	getPostResponse := new(GetPostcodeResponse)
	if err := json.Unmarshal(data, getPostResponse); err != nil {
		// Log the response
		log.Printf("%s", resp.Body)
		return nil, err
	}

	// If the status code was not 200, return error with the message from response
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(getPostResponse.Message)
	}

	return getPostResponse.Result, nil
}
