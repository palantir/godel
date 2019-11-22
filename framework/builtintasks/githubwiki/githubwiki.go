// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package githubwiki

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
	"github.com/pkg/errors"
)

type Params struct {
	DocsDir        string
	Repo           string
	AuthorName     string
	AuthorEmail    string
	CommitterName  string
	CommitterEmail string
	Msg            string
}

type GitTemplateParams struct {
	CommitID   string
	CommitTime time.Time
}

type commitUserParam struct {
	authorName     string
	authorEmail    string
	committerName  string
	committerEmail string
}

type userParam struct {
	envVar string
	format string
}

var (
	authorNameParam     = userParam{envVar: "GIT_AUTHOR_NAME", format: "an"}
	authorEmailParam    = userParam{envVar: "GIT_AUTHOR_EMAIL", format: "ae"}
	committerNameParam  = userParam{envVar: "GIT_COMMITTER_NAME", format: "cn"}
	committerEmailParam = userParam{envVar: "GIT_COMMITTER_EMAIL", format: "ce"}
)

type git string

func (g git) exec(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = string(g)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v failed\nError: %v\nOutput: %s", cmd.Args, err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

func (g git) clone(p Params) error {
	_, err := g.exec("clone", p.Repo, string(g))
	return err
}

func (g git) commitAll(msg string, p commitUserParam) error {
	// add all files
	if _, err := g.exec("add", "."); err != nil {
		return err
	}

	// commit to current branch
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Dir = string(g)
	env := os.Environ()
	// set environment variables for author and committer
	env = append(env, fmt.Sprintf("%v=%v", authorNameParam.envVar, p.authorName))
	env = append(env, fmt.Sprintf("%v=%v", authorEmailParam.envVar, p.authorEmail))
	env = append(env, fmt.Sprintf("%v=%v", committerNameParam.envVar, p.committerName))
	env = append(env, fmt.Sprintf("%v=%v", committerEmailParam.envVar, p.committerEmail))
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "%v failed: %s", cmd.Args, string(output))
	}

	return nil
}

func (g git) commitID(branch string) (string, error) {
	return g.exec("rev-parse", branch)
}

func (g git) commitTime(branch string) (time.Time, error) {
	output, err := g.exec("show", "-s", "--format=%ct", branch)
	if err != nil {
		return time.Time{}, err
	}
	unixTime, err := strconv.ParseInt(output, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to parse %s as an int64", output)
	}
	return time.Unix(unixTime, 0), nil
}

func (g git) push() error {
	_, err := g.exec("push", "origin", "HEAD")
	return err
}

func (g git) valueOr(v string, p userParam) (string, error) {
	if v != "" {
		return v, nil
	}
	output, err := g.exec("--no-pager", "show", "-s", fmt.Sprintf("--format=%%%v", p.format), "HEAD")
	if err != nil {
		return "", err
	}
	return output, nil
}

func createGitTemplateParams(docsDir string) (GitTemplateParams, error) {
	g := git(docsDir)

	commitID, err := g.commitID("HEAD")
	if err != nil {
		return GitTemplateParams{}, err
	}

	commitTime, err := g.commitTime("HEAD")
	if err != nil {
		return GitTemplateParams{}, err
	}

	return GitTemplateParams{
		CommitID:   commitID,
		CommitTime: commitTime,
	}, nil
}

func SyncGitHubWiki(p Params, stdout io.Writer) error {
	if err := layout.VerifyDirExists(p.DocsDir); err != nil {
		return errors.Wrapf(err, "Docs directory %s does not exist", p.DocsDir)
	}

	// apply templating to commit message
	msg := p.Msg
	if gitTemplateParams, err := createGitTemplateParams(p.DocsDir); err != nil {
		_, _ = fmt.Fprintf(stdout, "Failed to determine Git properties of documents directory %s: %v.\n", p.DocsDir, err)
		_, _ = fmt.Fprintln(stdout, "Continuing with templating disabled. To fix this issue, ensure that the directory is in a Git repository.")
	} else if t, err := template.New("message").Parse(p.Msg); err != nil {
		_, _ = fmt.Fprintf(stdout, "Failed to parse message %s as a template: %v. Using message as a string literal instead.\n", p.Msg, err)
	} else {
		buf := &bytes.Buffer{}
		if err := t.Execute(buf, gitTemplateParams); err != nil {
			_, _ = fmt.Fprintf(stdout, "Failed to execute template %s: %v. Using message as a string literal instead.\n", p.Msg, err)
		} else {
			msg = buf.String()
		}
	}

	tmpCloneDir, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	if err != nil {
		return errors.Wrapf(err, "Failed to create temporary directory")
	}

	g := git(tmpCloneDir)

	if err := g.clone(p); err != nil {
		return err
	}

	// update contents of cloned directory to match docs directory (ignoring .git directory)
	if modified, err := layout.SyncDir(p.DocsDir, tmpCloneDir, []string{".git"}); err != nil {
		return errors.Wrapf(err, "failed to sync contents of repo %s to docs directory %s", tmpCloneDir, p.DocsDir)
	} else if !modified {
		// nothing to do if cloned repo is identical to docs directory
		return nil
	}

	authorName, err := g.valueOr(p.AuthorName, authorNameParam)
	if err != nil {
		return errors.Wrapf(err, "failed to get authorName")
	}
	authorEmail, err := g.valueOr(p.AuthorEmail, authorEmailParam)
	if err != nil {
		return errors.Wrapf(err, "failed to get authorEmail")
	}
	committerName, err := g.valueOr(p.CommitterName, committerNameParam)
	if err != nil {
		return errors.Wrapf(err, "failed to get committerName")
	}
	committerEmail, err := g.valueOr(p.CommitterEmail, committerEmailParam)
	if err != nil {
		return errors.Wrapf(err, "failed to get committerEmail")
	}

	if err := g.commitAll(msg, commitUserParam{
		authorName:     authorName,
		authorEmail:    authorEmail,
		committerName:  committerName,
		committerEmail: committerEmail,
	}); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(stdout, "Pushing content of %s to %s...\n", p.DocsDir, p.Repo)
	return g.push()
}
