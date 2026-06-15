/*
MIT License

Copyright (c) 2026 gounix

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package jsonreq

import (
	"strings"
        "net/http"
        "encoding/json"
        "io"
	"strconv"
        "github-exporter/logger"
)

type RateLimitT struct {
	Limit int64
	Remaining int64
}

var limitStats RateLimitT

func GetRateLimit() RateLimitT {
	return limitStats
}

func ratelimit(header http.Header) {

        // 60 per hour for non authenticated calls, 5000 for authenticated
	// if breached failure code is 429 or 403
	str := header.Get("x-ratelimit-limit")
	value, err := strconv.Atoi(str)
	if err == nil {
		limitStats.Limit = int64(value)
	} else {
		logger.Error("jsonreq.ratelimit atoi x-ratelimit-limit", "err", err)
	}
	str = header.Get("x-ratelimit-remaining")
	value, err = strconv.Atoi(str)
	if err == nil {
		limitStats.Remaining = int64(value)
	} else {
		logger.Error("jsonreq.ratelimit atoi x-ratelimit-remaining", "err", err)
	}
}

func GetJsonResp(url string, token string, accept string, dat any) error {

	client := &http.Client{ }
        req, err := http.NewRequest("GET", url, nil)
	if accept != "" {
		req.Header.Add("accept", accept)
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer " + token)
	}

        resp, err := client.Do(req)
        if err != nil {
                logger.Error("jsonreq.getJsonResp", "client.do error", err)
                return err
        }

        defer resp.Body.Close()
	ratelimit(resp.Header)
        if resp.StatusCode != 200 {
                logger.Error("jsonreq.getJsonResp", "status", resp.Status)
                return err
        }

        body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("jsonreq.getJsonResp", "io.ReadAll error", err)
		return err
        }

        return json.Unmarshal(body, dat)
}


func getNextUrlFromHeader(header http.Header) string {
	// link_header := url ; relation [ , url ; relation ]
	// url         := <https://...>
	// relation    := rel="direction"
	// direction   := next | prev | last | first

	link_header := header.Get("link")
	if link_header == "" {
		return ""
	}

	// split on ,
	links := strings.Split(link_header, ",")
	for _, link := range links {

		link_parts := strings.Split(link, ";")
		encapsulated_url := link_parts[0]

		str_parts := strings.Split(encapsulated_url, "<")
		url2 := str_parts[1]
	        str_parts2 := strings.Split(url2, ">")
		url := str_parts2[0]

		relation := link_parts[1]

		rel_parts := strings.Split(relation, "\"")
		direction := rel_parts[1]

		if direction == "next" {
			return url
		}

	}
	return ""
}

func GetJsonRespPaginated(url string, token string, accept string, dat any) error {
	var full_body []map[string]interface{} 

	client := &http.Client{ }
	next_url := url
	for {
		var data []map[string]interface{} 

		logger.Info("GetJsonRespPaginated", "url", next_url)
		req, err := http.NewRequest("GET", next_url, nil)
		if accept != "" {
			req.Header.Add("accept", accept)
		}

		if token != "" {
			req.Header.Add("Authorization", "Bearer " + token)
		}

		resp, err := client.Do(req)
		if err != nil {
			logger.Error("jsonreq.GetJsonRespPaginated", "client.do error", err)
			return err
		}

		defer resp.Body.Close()
		ratelimit(resp.Header)
		if resp.StatusCode != 200 {
			logger.Error("jsonreq.GetJsonRespPaginated", "status", resp.Status)
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("jsonreq.GetJsonRespPaginated", "io.ReadAll error", err)
			return err
		}
		err = json.Unmarshal(body,&data)
		if err != nil {
			logger.Error("jsonreq.GetJsonRespPaginated", "json err", err)
		}

		full_body = append(full_body, data...)
		next_url = getNextUrlFromHeader(resp.Header)
		if next_url == "" {
			break
		}
	}
	jsonBody, err := json.Marshal(full_body)
	if err != nil {
		logger.Error("jsonreq.GetJsonRespPaginated", "json.Marshall err", err)
	}
        return json.Unmarshal(jsonBody, dat)
}
