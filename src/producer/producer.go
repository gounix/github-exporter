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
	"fmt"
	"github-exporter/data"
	"github-exporter/environ"
	"github-exporter/github"
	"github-exporter/jsonreq"
	"log/slog"
	"time"
)

const (
        trafficStatUrlPattern      = "https://api.github.com/repos/%s/traffic/%s" // project + clones or views
	repoStatUrlPattern         = "https://api.github.com/repos/%s"
	pullStatUrlPattern         = "https://api.github.com/repos/%s/pulls?state=all;per_page=250"
	issueStatUrlPattern        = "https://api.github.com/repos/%s/issues?state=all;per_page=250"
	tagStatUrlPattern          = "https://api.github.com/repos/%s/tags"
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
	TagT struct {
		Name string `json:"name"`
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
        slog.Info("producer.GetContributorsStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                slog.Error("producer.GetContributorsStats", "jsonreq.GetJsonResp", err)
                return []data.GHContributorStatsT{}, err
	}

	slog.Info("producer.GetContributorsStats", "project", project, "#contributors", len(raw))
	return raw, nil
}

func getCommitStats(project string) (int64, error) {
	var raw []CommitT

        url := fmt.Sprintf(commitStatUrlPattern, project)
        slog.Info("producer.getCommitStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                slog.Error("producer.getCommitStats", "jsonreq.GetJsonResp", err)
                return 0, err
	}

	slog.Info("producer.getCommitStats", "project", project, "#commits", len(raw))
	return int64(len(raw)), nil
}

func getTagStats(project string) (int64, error) {
	var raw []TagT

        url := fmt.Sprintf(tagStatUrlPattern, project)
        slog.Info("producer.getTagStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                slog.Error("producer.getTagStats", "jsonreq.GetJsonResp", err)
                return 0, err
	}

	slog.Info("producer.getTagStats", "project", project, "#tags", len(raw))
	return int64(len(raw)), nil
}

func getBranchStats(project string) (int64, error) {
	var raw []BranchT

        url := fmt.Sprintf(branchStatUrlPattern, project)
        slog.Info("producer.getBranchStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                slog.Error("producer.getBranchStats", "jsonreq.GetJsonResp", err)
                return 0, err
	}

	slog.Info("producer.getBranchStats", "project", project, "#branches", len(raw))
	return int64(len(raw)), nil
}

func getPullStats(project string) (data.GHPullStats, error) {
	var dat = data.GHPullStats{}
	var raw []PullsT

        url := fmt.Sprintf(pullStatUrlPattern, project)
        slog.Info("producer.getPullStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                slog.Error("producer.getPullStats", "jsonreq.GetJsonResp", err)
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

	slog.Info("producer.getPullStats", "project", project, "numOpen", dat.NumOpen, "numClosed", dat.NumClosed, "numUnassigned", dat.NumUnassigned)
	return dat, nil
}

func getIssueStats(project string) (data.GHIssueStats, error) {
	var dat = data.GHIssueStats{}
	var raw []IssuesT

        url := fmt.Sprintf(issueStatUrlPattern, project)
        slog.Info("producer.getIssueStats", "url", url)

	if err := jsonreq.GetJsonRespPaginated(url, environ.Env.Token, "application/vnd.github+json", &raw); err != nil {
                slog.Error("producer.getIssueStats", "jsonreq.GetJsonResp", err)
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

	slog.Info("producer.getIssueStats", "project", project, "numOpen", dat.NumOpen, "numClosed", dat.NumClosed, "numUnassigned", dat.NumUnassigned)
	return dat, nil
}

func getRepoStats(project string) (data.GHRepoStats, error) {
	var dat data.GHRepoStats

        url := fmt.Sprintf(repoStatUrlPattern, project)
        slog.Info("producer.getRepoStats", "url", url)

	if err := jsonreq.GetJsonResp(url, environ.Env.Token, "application/vnd.github+json", &dat); err != nil {
                slog.Error("producer.getRepoStats", "jsonreq.GetJsonResp", err)
                return data.GHRepoStats{}, err
	}

	slog.Info("producer.getRepoStats", "project", project, "StargazersCount", dat.StargazersCount, "WatchersCount", dat.WatchersCount, "ForksCount", dat.ForksCount, "OpenIssuesCount", dat.OpenIssuesCount, "NetworkCount", dat.NetworkCount, "SubscribersCount", dat.SubscribersCount)
	return dat, nil
}

func getTrafficStats(project string, stat string) (data.GHTrafficStats, error) {
	var dat data.GHTrafficStats

        url := fmt.Sprintf(trafficStatUrlPattern, project, stat)
        slog.Info("producer.getTrafficStats", "url", url)

	if err := jsonreq.GetJsonResp(url, environ.Env.Token, "application/vnd.github+json", &dat); err != nil {
                slog.Error("producer.getTrafficStats", "jsonreq.GetJsonResp", err)
                return data.GHTrafficStats{}, err
	}

	slog.Info("producer.getTrafficStats", "project", project, "stat", stat, "count", dat.Count, "uniques", dat.Uniques)
	return dat, nil
}



func repoLoop(repos []string) {
	var err error

	for _, project := range repos {
		var stats data.ProjectT

		stats.Project = project
		stats.Clones, err = getTrafficStats(project, "clones")
		if err != nil {
			slog.Error("producer.repoLoop clones", "project", project, "err", err)
		}
		stats.Views, err = getTrafficStats(project, "views")
		if err != nil {
			slog.Error("producer.repoLoop views", "project", project, "err", err)
		}
		stats.RepoStats, err = getRepoStats(project)
		if err != nil {
			slog.Error("producer.repoLoop repostats", "project", project, "err", err)
		}

		stats.PullStats, err = getPullStats(project)
		if err != nil {
			slog.Error("producer.repoLoop pullstats", "project", project, "err", err)
		}

		stats.IssueStats, err = getIssueStats(project)
		if err != nil {
			slog.Error("producer.repoLoop issuestats", "project", project, "err", err)
		}

		stats.TagStats, err = getTagStats(project)
		if err != nil {
			slog.Error("producer.repoLoop tagstats", "project", project, "err", err)
		}

		stats.BranchStats, err = getBranchStats(project)
		if err != nil {
			slog.Error("producer.repoLoop branchstats", "project", project, "err", err)
		}

		stats.CommitStats, err = getCommitStats(project)
		if err != nil {
			slog.Error("producer.repoLoop commitstats", "project", project, "err", err)
		}

		stats.ContributorStats, err = GetContributorsStats(project)
		if err != nil {
			slog.Error("producer.repoLoop contributorsStats", "project", project, "err", err)
		}

		limitStats := jsonreq.GetRateLimit()
		slog.Info("producer.repoLoop", "x-ratelimit-limit", limitStats.Limit, "x-ratelimit-remaining", limitStats.Remaining)

		data.Put(stats, limitStats)
	}
}

func Refresh() {
	for {
		repos, err := github.GetRepos(environ.Env.GithubUser)
		if err != nil {
			slog.Error("producer/Refresh", "github.GetRepos", err)
		} else {
			data.Initialize(repos)
			repoLoop(repos)
		}

		slog.Info("producer.Refresh", "sleeping", environ.Env.RefreshSeconds)
		time.Sleep(time.Duration(environ.Env.RefreshSeconds) * time.Second)
	}
}
