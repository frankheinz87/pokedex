package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type Area struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Response struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Area  `json:"results"`
}

func (c *Client) ListLocations(url string) (Response, error) {
	// 1. cache first
	if data, ok := c.cache.Get(url); ok {
		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			return Response{}, err
		}
		return resp, nil
	}

	// 2. cache miss → HTTP
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return Response{}, err
	}

	// 3. store in cache
	c.cache.Add(url, data)

	return resp, nil
}
