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

package producer

import (
	"github-exporter/data"
	"github-exporter/logger"
	"github-exporter/environ"
	"github-exporter/jsonreq"
	"time"
	"fmt"
)

const (
        statUrlPattern    = "https://api.github.com/repos/%s/traffic/%s" // project + clones or views
)


func getStat(project string, stat string) (data.GHStats, error) {
	var dat data.GHStats

        url := fmt.Sprintf(statUrlPattern, project, stat)
        logger.Info("producer.getStat", "url", url)

	if err := jsonreq.GetJsonResp(url, environ.Env.Token, "application/vnd.github+json", &dat); err != nil {
                logger.Error("producer.getStat", "jsonreq.GetJsonResp", err)
                return data.GHStats{}, err
	}

	return dat, nil
}



func Put(repos []string, interval int64) {
	data.Initialize(repos)
	for {
		for _, entry := range repos {
			clones, cl_err := getStat(entry, "clones")
			views, v_err := getStat(entry, "views")
			if cl_err == nil && v_err == nil{
				data.Put(entry, clones, views)
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
