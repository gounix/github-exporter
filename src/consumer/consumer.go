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

package consumer

import (
	"fmt"
	"net/http"
	"github-exporter/data"
	"github-exporter/logger"
	"github-exporter/jsonreq"
)

const promHeader = `# HELP github_clones number of clones of github project
# TYPE github_clones gauge
# HELP github_repo_stats count of stat
# TYPE github_repo_stats gauge
# HELP github_views number of views of github project
# TYPE github_views gauge`

var interval int64

func logRequest(r *http.Request) {
	logger.Info("consumer.logRequest", "host", r.Host, "method", r.Method, "url", r.URL.Path, "user agent", r.UserAgent())
}

func sendSingleLine(w http.ResponseWriter, project string, stat string, value int64) {

	str := fmt.Sprintf("github_repo_stats{project=\"%s\",stat=\"%s\"} %d\n", project, stat, value)
	fmt.Fprintf(w, str)
	logger.Info("consumer.sendSingleLine", "reply", str)
}

func sendPromLines(w http.ResponseWriter, limitStats jsonreq.RateLimitT, projects []data.ProjectT) {

	for _, entry := range projects {
		str := fmt.Sprintf("github_exporter_stats{stat=\"ratelimit-limit\"} %d\n", limitStats.Limit)
		fmt.Fprintf(w, str)
		logger.Info("consumer.sendPromLines", "reply", str)
		str = fmt.Sprintf("github_exporter_stats{stat=\"ratelimit-remaining\"} %d\n", limitStats.Remaining)
		fmt.Fprintf(w, str)
		logger.Info("consumer.sendPromLines", "reply", str)

		str = fmt.Sprintf("github_clones{project=\"%s\",unique=\"true\"} %d\n", entry.Project, entry.Clones.Uniques)
		fmt.Fprintf(w, str)
		logger.Info("consumer.sendPromLines", "reply", str)
		str = fmt.Sprintf("github_clones{project=\"%s\",unique=\"false\"} %d\n", entry.Project, entry.Clones.Count)
		fmt.Fprintf(w, str)
		logger.Info("consumer.sendPromLines", "reply", str)

		str = fmt.Sprintf("github_views{project=\"%s\",unique=\"true\"} %d\n", entry.Project, entry.Views.Uniques)
		fmt.Fprintf(w, str)
		logger.Info("consumer.sendPromLines", "reply", str)
		str = fmt.Sprintf("github_views{project=\"%s\",unique=\"false\"} %d\n", entry.Project, entry.Views.Count)
		fmt.Fprintf(w, str)
		logger.Info("consumer.sendPromLines", "reply", str)

		sendSingleLine(w, entry.Project, "stargazers_count", entry.RepoStats.StargazersCount)
		sendSingleLine(w, entry.Project, "watchers_count", entry.RepoStats.WatchersCount)
		sendSingleLine(w, entry.Project, "forks_count", entry.RepoStats.ForksCount)
		sendSingleLine(w, entry.Project, "network_count", entry.RepoStats.NetworkCount)
		sendSingleLine(w, entry.Project, "subscribers_count", entry.RepoStats.SubscribersCount)
		sendSingleLine(w, entry.Project, "open_pull_requests", entry.PullStats.NumOpen)
		sendSingleLine(w, entry.Project, "closed_pull_requests", entry.PullStats.NumClosed)
		sendSingleLine(w, entry.Project, "unassigned_pull_requests", entry.PullStats.NumUnassigned)
		sendSingleLine(w, entry.Project, "open_issues", entry.IssueStats.NumOpen)
		sendSingleLine(w, entry.Project, "closed_issues", entry.IssueStats.NumClosed)
		sendSingleLine(w, entry.Project, "unassigned_issues", entry.IssueStats.NumUnassigned)
		sendSingleLine(w, entry.Project, "branches", entry.BranchStats)
		sendSingleLine(w, entry.Project, "commits", entry.CommitStats)

		for _, contributor := range entry.ContributorStats {
			str = fmt.Sprintf("github_repo_stats{project=\"%s\",stat=\"contributors\",contributor=\"%s\"} %d\n", entry.Project, contributor.Login, contributor.Contributions)
			fmt.Fprintf(w, str)
			logger.Info("consumer.sendPromLines", "reply", str)
		}
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {

	logRequest(r)
	limitStats, projects, err := data.Get()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		fmt.Fprintln(w, promHeader)
		sendPromLines(w, limitStats, projects)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if data.Alive(interval) {
		fmt.Fprintf(w, "OK")
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Get(port int64, intrvl int64) {
	interval = intrvl
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/health", healthHandler)
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}

