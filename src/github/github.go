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

package github

import (
	"fmt"
	"github-exporter/jsonreq"
	"github-exporter/logger"
)

const getReposUrlPattern = "https://api.github.com/users/%s/repos"

type ReposT struct {
	FullName string `json:"full_name"`
}

func GetRepos(user string) ([]string, error) {
	var dat []ReposT
	var lst []string

	url := fmt.Sprintf(getReposUrlPattern, user)
	logger.Info("github.GetRepos", "url", url)

	if err := jsonreq.GetJsonResp(url, "", "application/vnd.github+json", &dat); err != nil {
                logger.Error("github.GetRepos", "jsonreq.GetJsonResp", err)
                return []string{}, err
        }

	for _, entry := range dat {
		logger.Info("github-exporter/GetRepos", "repo", entry.FullName)
		lst = append(lst, entry.FullName)
	}
        return lst, nil

}
