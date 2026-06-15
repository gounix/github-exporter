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
        trafficStatUrlPattern      = "https://api.github.com/repos/%s/traffic/%s" // project + clones or views
	repoStatUrlPattern         = "https://api.github.com/repos/%s"
	pullStatUrlPattern         = "https://api.github.com/repos/%s/pulls?state=all;per_page=250"
	issueStatUrlPattern        = "https://api.github.com/repos/%s/issues?state=all;per_page=250"
	branchStatUrlPattern       = "https://api.github.com/repos/%s/branches"
	commitStatUrlPattern       = "https://api.github.com/repos/%s/commits?per_page=250"
	contributorsStatUrlPattern = "https://api.github.com/repos/%s/contributors?per_page=250"
)

type (
	CommitterT struct {
		Login string `json:"login"`
	}
	CommitT struct {
		Committer CommitterT `json:"committer"`
	}
	BranchT struct {
		Name string `json:"name"`
	}
	AssigneeT struct {
		Login string `json:"login"`
	}
	PullsT struct {
		State string `json:"state"`
		Assignees []AssigneeT `json:"assignees"`
	}
	PullRequestT struct {
		Url string `json:"url"`
	}
	IssuesT struct {
		State string `json:"state"`
		Title string `json:"title"`
		Id int64 `json:"id"`
		PullRequest PullRequestT `json:"pull_request"`
		Assignees []AssigneeT `json:"assignees"`
	}
)

func GetContributorsStats(project string) ([]data.GHContributorStatsT, error) {
	var raw []data.GHContributorStatsT

        url := fmt.Sprintf(contributorsStatUrlPattern, project)
        logger.Info("producer.GetContributorsStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                logger.Error("producer.GetContributorsStats", "jsonreq.GetJsonResp", err)
                return []data.GHContributorStatsT{}, err
	}

	return raw, nil
}

func getCommitStats(project string) (int64, error) {
	var raw []CommitT

        url := fmt.Sprintf(commitStatUrlPattern, project)
        logger.Info("producer.getCommitStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                logger.Error("producer.getCommitStats", "jsonreq.GetJsonResp", err)
                return 0, err
	}

	return int64(len(raw)), nil
}

func getBranchStats(project string) (int64, error) {
	var raw []BranchT

        url := fmt.Sprintf(branchStatUrlPattern, project)
        logger.Info("producer.getBranchStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                logger.Error("producer.getBranchStats", "jsonreq.GetJsonResp", err)
                return 0, err
	}

	return int64(len(raw)), nil
}

func getPullStats(project string) (data.GHPullStats, error) {
	var dat = data.GHPullStats{}
	var raw []PullsT

        url := fmt.Sprintf(pullStatUrlPattern, project)
        logger.Info("producer.getPullStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                logger.Error("producer.getPullStats", "jsonreq.GetJsonResp", err)
                return data.GHPullStats{}, err
	}
	for _, entry := range raw {
		if entry.State == "open" {
			dat.NumOpen++
			if len(entry.Assignees) == 0 {
				dat.NumUnassigned++
			}
		} else {
			dat.NumClosed++
		}
	}

	return dat, nil
}

func getIssueStats(project string) (data.GHIssueStats, error) {
	var dat = data.GHIssueStats{}
	var raw []IssuesT

        url := fmt.Sprintf(issueStatUrlPattern, project)
        logger.Info("producer.getIssueStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                logger.Error("producer.getIssueStats", "jsonreq.GetJsonResp", err)
                return data.GHIssueStats{}, err
	}
	for _, entry := range raw {
		// check if it is not a pull request
		if entry.PullRequest.Url == "" {
			if entry.State == "open" {
				dat.NumOpen++
				if len(entry.Assignees) == 0 {
					dat.NumUnassigned++
				}
			} else {
				dat.NumClosed++
			}
		}
	}

	return dat, nil
}

func getRepoStats(project string) (data.GHRepoStats, error) {
	var dat data.GHRepoStats

        url := fmt.Sprintf(repoStatUrlPattern, project)
        logger.Info("producer.getRepoStats", "url", url)

	if err := jsonreq.GetJsonResp(url, environ.Env.Token, "application/vnd.github+json", &dat); err != nil {
                logger.Error("producer.getRepoStats", "jsonreq.GetJsonResp", err)
                return data.GHRepoStats{}, err
	}

	return dat, nil
}

func getTrafficStat(project string, stat string) (data.GHTrafficStats, error) {
	var dat data.GHTrafficStats

        url := fmt.Sprintf(trafficStatUrlPattern, project, stat)
        logger.Info("producer.getTrafficStat", "url", url)

	if err := jsonreq.GetJsonResp(url, environ.Env.Token, "application/vnd.github+json", &dat); err != nil {
                logger.Error("producer.getTrafficStat", "jsonreq.GetJsonResp", err)
                return data.GHTrafficStats{}, err
	}

	return dat, nil
}



func Put(repos []string, interval int64) {
	data.Initialize(repos)
	for {
		for _, entry := range repos {
			clones, err := getTrafficStat(entry, "clones")
			if err != nil {
				logger.Error("producer.Put clones", "project", entry, "err", err)
			}
			views, err := getTrafficStat(entry, "views")
			if err != nil {
				logger.Error("producer.Put views", "project", entry, "err", err)
			}
			repoStats, err := getRepoStats(entry)
			if err != nil {
				logger.Error("producer.Put repostats", "project", entry, "err", err)
			}

			pullStats, err := getPullStats(entry)
			if err != nil {
				logger.Error("producer.Put pullstats", "project", entry, "err", err)
			}
			logger.Info("producer.Put", "project pull stats", entry, "numOpen", pullStats.NumOpen, "numClosed", pullStats.NumClosed, "numUnassigned", pullStats.NumUnassigned)

			issueStats, err := getIssueStats(entry)
			if err != nil {
				logger.Error("producer.Put issuestats", "project", entry, "err", err)
			}
			logger.Info("producer.Put", "project issue stats", entry, "numOpen", issueStats.NumOpen, "numClosed", issueStats.NumClosed, "numUnassigned", issueStats.NumUnassigned)

			branchStats, err := getBranchStats(entry)
			if err != nil {
				logger.Error("producer.Put branchstats", "project", entry, "err", err)
			}
			logger.Info("producer.Put", "branch stats", entry, "#branches", branchStats)

			commitStats, err := getCommitStats(entry)
			if err != nil {
				logger.Error("producer.Put commitstats", "project", entry, "err", err)
			}
			logger.Info("producer.Put", "commit stats", entry, "#commits", commitStats)

			contributorStats, err := GetContributorsStats(entry)
			if err != nil {
				logger.Error("producer.Put contributorsStats", "project", entry, "err", err)
			}
			logger.Info("producer.Put", "contributor stats", entry, "#contributors", len(contributorStats))

			limitStats := jsonreq.GetRateLimit()
			logger.Info("producer.Put", "x-ratelimit-limit", limitStats.Limit, "x-ratelimit-remaining", limitStats.Remaining)

			data.Put(entry, clones, views, repoStats, pullStats, issueStats, branchStats, commitStats, contributorStats, limitStats)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
