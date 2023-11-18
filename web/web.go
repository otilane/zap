package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Method string

const (
	Get    Method = http.MethodGet
	Post   Method = http.MethodPost
	Put    Method = http.MethodPut
	Delete Method = http.MethodDelete
)

// Helper function for creeating a HTTP request and returning the body of the result
func FetchBody(method Method, apiURL string, bodyReader io.Reader, headers map[string]string) ([]byte, error) {

	// Create a new GET request with the provided URL
	req, err := http.NewRequest(fmt.Sprint(method), apiURL, bodyReader)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making %v request:", method)
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}

// Helper function for setting parameters in the url and encoding them
func SetParams(apiURL string, params map[string]string) string {
	p := url.Values{}

	// Add parameters
	for key, value := range params {
		p.Add(key, value)
	}
	return apiURL + p.Encode()
}
