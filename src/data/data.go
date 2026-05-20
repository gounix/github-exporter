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

package data

import (
	"sync"
	"time"
	"errors"
	"github-exporter/logger"
)

type (
	GHStats struct {
                Count   int64 `json:"count"`
                Uniques int64 `json:"uniques"`
        }
	ProjectT struct {
		Project     string
		Clones      GHStats
		Views       GHStats
		Initialized bool
	}
	dataT struct {
		mu          sync.Mutex
		projects    []ProjectT
		timestamp   time.Time
		initialized bool
	}
)

var data = dataT{ initialized: false }

func Initialize(projects []string) {
	for _, entry := range projects {
		var project_data ProjectT
		project_data.Project = entry
		project_data.Initialized = false
		data.projects = append(data.projects, project_data)
	}
}

func Put(project string, clones GHStats, views GHStats) {
	data.mu.Lock()
	defer data.mu.Unlock()

	for nr, _ := range data.projects {
		if data.projects[nr].Project == project {
			data.projects[nr].Clones = clones
			data.projects[nr].Views = views
			data.projects[nr].Initialized = true

			data.timestamp = time.Now()
			data.initialized = true

			logger.Info("data.Put", "project", data.projects[nr].Project, "clones.count", data.projects[nr].Clones.Count, "clones.uniques", data.projects[nr].Clones.Uniques, "views.count", data.projects[nr].Views.Count, "views.uniques", data.projects[nr].Views.Uniques)
		}
	}
}

func Get() ([]ProjectT, error) {
	data.mu.Lock()
	defer data.mu.Unlock()

	if !data.initialized {
		logger.Info("data.Get all not initialized")
		return []ProjectT{}, errors.New("not initialized")
	}
	for _, entry := range data.projects {
		if !entry.Initialized {
			logger.Info("data.Get not initialized", "entry", entry.Project, "clones.count", entry.Clones.Count)
			return []ProjectT{}, errors.New("not initialized")
		}
	}
	// all is OK now
	logger.Info("data.Get", "first project", data.projects[0].Project, "nr projects", len(data.projects))
	return data.projects, nil
}

func Alive(interval int64) bool {
	now := time.Now()
	diff := now.Sub(data.timestamp)

	// consider the producer dead after it missed 2 intervals
	isOK := diff.Seconds() < float64(2 * interval)
	logger.Info("data.Alive", "age", diff.Seconds(), "OK", isOK)

	return isOK
}

