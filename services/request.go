package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// User agent used when communicating with the Paystack API.
	userAgent = "Mozilla/5.0 (Unknown; Linux) AppleWebKit/538.1 (KHTML, like Gecko) Chrome/v1.0.0 Safari/538.1"
)

func GetResource(ctx context.Context, url string, res interface{}, token string) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	return doReq(req, res, token)
}

func doReq(req *http.Request, res interface{}, token string) error {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error processing request - %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected http status: %s", resp.Status)
		return fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	err = parseAPIResponse(resp, res)
	if err != nil {
		return err
	}

	return nil
}

func parseAPIResponse(resp *http.Response, response interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(&response)

	fmt.Println(response)

	if err != nil {
		return fmt.Errorf("error while unmarshalling the response bytes %+v ", err)
	}

	return nil
}
