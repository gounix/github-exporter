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
	"math"
	"sync"
	"time"
	"errors"
	"github-exporter/logger"
	"github-exporter/jsonreq"
)

type (
	GHRepoStats struct {
                StargazersCount  int64 `json:"stargazers_count"`
                WatchersCount    int64 `json:"watchers_count"`
                ForksCount       int64 `json:"forks_count"`
                OpenIssuesCount  int64 `json:"open_issues_count"`
                NetworkCount     int64 `json:"network_count"`
                SubscribersCount int64 `json:"subscribers_count"`
        }
	GHTrafficStats struct {
                Count   int64 `json:"count"`
                Uniques int64 `json:"uniques"`
        }
	GHPullStats struct {
		NumOpen int64
		NumClosed int64
		NumUnassigned int64
	}
	GHIssueStats struct {
		NumOpen int64
		NumClosed int64
		NumUnassigned int64
	}
	GHContributorStatsT struct {
                Login string `json:"login"`
                Contributions int64 `json:""contributions""`
        }

	ProjectT struct {
		Project     string
		Clones      GHTrafficStats
		Views       GHTrafficStats
		RepoStats   GHRepoStats
		PullStats   GHPullStats
		IssueStats  GHIssueStats
		ContributorStats  []GHContributorStatsT
		BranchStats int64
		CommitStats int64
		Initialized bool
	}
	dataT struct {
		mu          sync.Mutex
		projects    []ProjectT
		timestamp   time.Time
		limitStats  jsonreq.RateLimitT
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

func Put(project string, clones GHTrafficStats, views GHTrafficStats, repoStats GHRepoStats, pullStats GHPullStats, issueStats GHIssueStats, branchStats int64, commitStats int64, contributorStats []GHContributorStatsT, limitStats jsonreq.RateLimitT) {
	data.mu.Lock()
	defer data.mu.Unlock()

	data.limitStats = limitStats
	for nr, _ := range data.projects {
		if data.projects[nr].Project == project {
			data.projects[nr].Clones = clones
			data.projects[nr].Views = views
			data.projects[nr].RepoStats = repoStats
			data.projects[nr].PullStats = pullStats
			data.projects[nr].IssueStats = issueStats
			data.projects[nr].BranchStats = branchStats
			data.projects[nr].CommitStats = commitStats
			data.projects[nr].ContributorStats = contributorStats
			data.projects[nr].Initialized = true

			data.timestamp = time.Now()
			data.initialized = true

			logger.Info("data.Put", "project", data.projects[nr].Project, "clones.count", data.projects[nr].Clones.Count, "clones.uniques", data.projects[nr].Clones.Uniques, "views.count", data.projects[nr].Views.Count, "views.uniques", data.projects[nr].Views.Uniques)
		}
	}
}

func Get() (jsonreq.RateLimitT, []ProjectT, error) {
	data.mu.Lock()
	defer data.mu.Unlock()

	if !data.initialized {
		logger.Info("data.Get all not initialized")
		return jsonreq.RateLimitT{}, []ProjectT{}, errors.New("not initialized")
	}
	for _, entry := range data.projects {
		if !entry.Initialized {
			logger.Info("data.Get not initialized", "entry", entry.Project, "clones.count", entry.Clones.Count)
			return jsonreq.RateLimitT{}, []ProjectT{}, errors.New("not initialized")
		}
	}
	// all is OK now
	logger.Info("data.Get", "first project", data.projects[0].Project, "nr projects", len(data.projects))
	return data.limitStats, data.projects, nil
}

func Alive(interval int64) bool {
	now := time.Now()
	diff := now.Sub(data.timestamp)

	// consider the producer dead after it missed 2 intervals
	isOK := diff.Seconds() < float64(2 * interval)
	logger.Info("data.Alive", "age(seconds)", math.Floor(diff.Seconds()), "OK", isOK)

	return isOK
}

