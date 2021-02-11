package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = ".instructure.com"

// CanvasClient stores data releveant to the operation of canvas api
type CanvasClient struct {
	Domain  string
	client  http.Client
	headers http.Header
}

// NewClient creates new client
func NewClient(domain string, authorizationToken string) *CanvasClient {
	c := CanvasClient{
		Domain:  domain,
		client:  *http.DefaultClient,
		headers: http.Header{},
	}

	c.headers.Add("authorization", "Bearer "+authorizationToken)

	return &c
}

// ClientURL returns a complete client URL
func (c *CanvasClient) ClientURL() string {
	return fmt.Sprintf("https://%s.instructure.com", c.Domain)
}

func (c *CanvasClient) getJSON(url string, target interface{}) error {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.Header = c.headers
	res, err := c.client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Status code is: %d", res.StatusCode)
	}
	defer res.Body.Close()
	// responseData, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Print(string(responseData))

	return json.NewDecoder(res.Body).Decode(&target)
}
