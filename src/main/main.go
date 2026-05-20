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

package main

import (
	"os"
	"github-exporter/producer"
	"github-exporter/consumer"
	"github-exporter/environ"
	"github-exporter/logger"
	"github-exporter/github"
)

func main() {
	if err := environ.Load(); err != nil {
                logger.Error("github-exporter/main", "environ.Load", err)
		os.Exit(1)
        }

	repos, err := github.GetRepos(environ.Env.GithubUser)
	if err != nil {
		logger.Error("github-exporter/main", "github.GetRepos", err)
		os.Exit(1)
	}

	go producer.Put(repos, environ.Env.RefreshSeconds)
	consumer.Get(environ.Env.PortNumber, environ.Env.RefreshSeconds)
}
