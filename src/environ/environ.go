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

package environ

import (
        "go-simpler.org/env"
	"github-exporter/logger"
	"os"
)

type EnvT struct {
        Token          string `env:"TOKEN,required"`
	RefreshSeconds int64  `env:"REFRESH_SECONDS" default:"300"`
        PortNumber     int64  `env:"PORT_NUMBER" default:"9900"`
        GithubUser     string `env:"GITHUB_USER,required"`
}

var Env EnvT

func Load() error {
	if err := env.Load(&Env, nil); err != nil {
                logger.Error("github-exporter/environ", "env.Load", err)
		os.Exit(1)
        }
	logger.Info("github-exporter.environ loaded environment", "TOKEN(truncated)", Env.Token[:5])
	logger.Info("github-exporter.environ loaded environment", "REFRESH_SECONDS", Env.RefreshSeconds)
	logger.Info("github-exporter.environ loaded environment", "PORT_NUMBER", Env.PortNumber)
	logger.Info("github-exporter.environ loaded environment", "GITHUB_USER", Env.GithubUser)
	return nil
}
