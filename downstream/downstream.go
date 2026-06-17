package downstream

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Target struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

type Result struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
	Body   string `json:"body"`
	Error  string `json:"error,omitempty"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func CallAll(targets []Target, payload []byte) []Result {
	results := make([]Result, 0, len(targets))

	for _, t := range targets {
		r := Result{URL: t.URL}

		method := t.Method
		if method == "" {
			method = "POST"
		}

		var bodyReader io.Reader
		if payload != nil {
			bodyReader = bytes.NewReader(payload)
		}

		req, err := http.NewRequest(method, t.URL, bodyReader)
		if err != nil {
			r.Error = fmt.Sprintf("create request: %v", err)
			results = append(results, r)
			continue
		}

		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			r.Error = fmt.Sprintf("execute: %v", err)
			results = append(results, r)
			continue
		}

		r.Status = resp.StatusCode
		respBody, _ := io.ReadAll(resp.Body)
		r.Body = string(respBody)
		resp.Body.Close()

		results = append(results, r)
	}

	return results
}
