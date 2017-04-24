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

package publish

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/cmd/dist"
	"github.com/palantir/godel/apps/distgo/config"
	"github.com/palantir/godel/apps/distgo/params"
)

const (
	urlFlagName            = "url"
	userFlagName           = "user"
	passwordFlagName       = "password"
	almanacURLFlagName     = "almanac-url"
	almanacIDFlagName      = "almanac-id"
	almanacSecretFlagName  = "almanac-secret"
	almanacReleaseFlagName = "almanac-release"
	repositoryFlagName     = "repository"
	subjectFlagName        = "subject"
	publishFlagName        = "publish"
	downloadsListFlagName  = "downloads-list"
	pathFlagName           = "path"
	failFastFlagName       = "fail-fast"
	gitHubOwnerFlagName    = "owner"
)

var (
	urlFlag = flag.StringFlag{
		Name:     urlFlagName,
		Usage:    "URL for a remote repository (such as https://repository.domain.com)",
		Required: true,
	}
	userFlag = flag.StringFlag{
		Name:     userFlagName,
		Usage:    "Username for authentication for repository",
		Required: true,
	}
	passwordFlag = flag.StringFlag{
		Name:     passwordFlagName,
		Usage:    "Password for repository",
		Required: true,
	}
	almanacURLFlag = flag.StringFlag{
		Name:     almanacURLFlagName,
		Usage:    "URL for an Almanac instance (such as https://almanac.domain.com)",
		Required: false,
	}
	almanacIDFlag = flag.StringFlag{
		Name:     almanacIDFlagName,
		Usage:    "Almanac access ID",
		Required: false,
	}
	almanacSecretFlag = flag.StringFlag{
		Name:     almanacSecretFlagName,
		Usage:    "Almanac secret",
		Required: false,
	}
	almanacReleaseFlag = flag.BoolFlag{
		Name:  almanacReleaseFlagName,
		Usage: "Perform an Almanac release after publish",
	}
	failFastFlag = flag.BoolFlag{
		Name:  failFastFlagName,
		Usage: "Fail immediately if the publish operation for an individual product fails",
	}
)

func Command() cli.Command {
	publishCmd := cli.Command{
		Name:  "publish",
		Usage: "Publish product distributions",
		Subcommands: []cli.Command{
			localType.createCommand(),
			artifactoryType.createCommand(),
			bintrayType.createCommand(),
			githubType.createCommand(),
		},
	}

	for i := range publishCmd.Subcommands {
		publishCmd.Subcommands[i].Flags = append(publishCmd.Subcommands[i].Flags, cmd.ProductsParam)
	}

	return publishCmd
}

type publisherType struct {
	name      string
	usage     string
	flags     []flag.Flag
	publisher func(ctx cli.Context) Publisher
}

func (p *publisherType) createCommand() cli.Command {
	return cli.Command{
		Name:  p.name,
		Usage: p.usage,
		Flags: p.flags,
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return publishAction(p.publisher(ctx), ctx.Slice(cmd.ProductsParamName), newAlmanacInfo(ctx), ctx.Bool(failFastFlagName), ctx.App.Stdout, wd)
		},
	}
}

var (
	localType = publisherType{
		name:  "local",
		usage: "Publish products to a local directory",
		flags: []flag.Flag{
			flag.StringFlag{
				Name:  pathFlagName,
				Usage: "Path to local publish root",
				Value: path.Join(os.Getenv("HOME"), ".m2", "repository"),
			},
			failFastFlag,
		},
		publisher: func(ctx cli.Context) Publisher {
			return &LocalPublishInfo{
				Path: ctx.String(pathFlagName),
			}
		},
	}
	artifactoryType = publisherType{
		name:  "artifactory",
		usage: "Publish products to an Artifactory repository",
		flags: remotePublishFlags(flag.StringFlag{
			Name:     repositoryFlagName,
			Usage:    "Repository that is the destination for the publish",
			Required: true,
		}),
		publisher: func(ctx cli.Context) Publisher {
			return &ArtifactoryConnectionInfo{
				BasicConnectionInfo: basicRemoteInfo(ctx),
				Repository:          ctx.String(repositoryFlagName),
			}
		},
	}
	bintrayType = publisherType{
		name:  "bintray",
		usage: "Publish products to a Bintray repository",
		flags: remotePublishFlags(
			flag.StringFlag{
				Name:     subjectFlagName,
				Usage:    "Subject that is the destination for the publish",
				Required: true,
			},
			flag.StringFlag{
				Name:     repositoryFlagName,
				Usage:    "Repository that is the destination for the publish",
				Required: true,
			},
			flag.BoolFlag{
				Name:  publishFlagName,
				Usage: "Publish uploaded content",
			},
			flag.BoolFlag{
				Name:  downloadsListFlagName,
				Usage: "Add uploaded artifact to downloads list for package",
			},
		),
		publisher: func(ctx cli.Context) Publisher {
			return &BintrayConnectionInfo{
				BasicConnectionInfo: basicRemoteInfo(ctx),
				Subject:             ctx.String(subjectFlagName),
				Repository:          ctx.String(repositoryFlagName),
				Release:             ctx.Bool(publishFlagName),
				DownloadsList:       ctx.Bool(downloadsListFlagName),
			}
		},
	}
	githubType = publisherType{
		name:  "github",
		usage: "Publish products to a GitHub repository",
		flags: remotePublishFlags(
			flag.StringFlag{
				Name:     repositoryFlagName,
				Usage:    "Repository that is the destination for the publish",
				Required: true,
			},
			flag.StringFlag{
				Name:  gitHubOwnerFlagName,
				Usage: "GitHub owner of the destination repository for the publish (if unspecified, user will be used)",
			},
		),
		publisher: func(ctx cli.Context) Publisher {
			return &GitHubConnectionInfo{
				APIURL:     basicRemoteInfo(ctx).URL,
				User:       basicRemoteInfo(ctx).Username,
				Token:      basicRemoteInfo(ctx).Password,
				Owner:      ctx.String(gitHubOwnerFlagName),
				Repository: ctx.String(repositoryFlagName),
			}
		},
	}
)

func remotePublishFlags(flags ...flag.Flag) []flag.Flag {
	remoteFlags := []flag.Flag{
		urlFlag,
		userFlag,
		passwordFlag,
	}
	remoteFlags = append(remoteFlags, flags...)
	remoteFlags = append(remoteFlags,
		failFastFlag,
		almanacURLFlag,
		almanacIDFlag,
		almanacSecretFlag,
		almanacReleaseFlag,
	)
	return remoteFlags
}

func basicRemoteInfo(ctx cli.Context) BasicConnectionInfo {
	return BasicConnectionInfo{
		URL:      ctx.String(urlFlagName),
		Username: ctx.String(userFlagName),
		Password: ctx.String(passwordFlagName),
	}
}

func newAlmanacInfo(ctx cli.Context) *AlmanacInfo {
	if !ctx.Has(almanacURLFlagName) || ctx.String(almanacURLFlagName) == "" {
		return nil
	}
	return &AlmanacInfo{
		URL:      ctx.String(almanacURLFlagName),
		AccessID: ctx.String(almanacIDFlagName),
		Secret:   ctx.String(almanacSecretFlagName),
		Release:  ctx.Bool(almanacReleaseFlagName),
	}
}

func publishAction(publisher Publisher, products []string, almanacInfo *AlmanacInfo, failFast bool, stdout io.Writer, wd string) error {
	cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
	if err != nil {
		return err
	}

	return build.RunBuildFunc(func(buildSpecWithDeps []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
		distsNotBuilt := DistsNotBuilt(buildSpecWithDeps)
		var specsToBuild []params.ProductBuildSpec
		for _, currSpecWithDeps := range distsNotBuilt {
			specsToBuild = append(specsToBuild, build.RequiresBuild(currSpecWithDeps, nil).Specs()...)
		}
		if len(specsToBuild) > 0 {
			if err := build.Run(specsToBuild, nil, build.DefaultContext(), stdout); err != nil {
				return errors.Wrapf(err, "failed to build products required for dist")
			}
		}

		if err := cmd.ProcessSerially(dist.Run)(distsNotBuilt, stdout); err != nil {
			return errors.Wrapf(err, "failed to build dists required for publish")
		}

		var processFunc cmd.ProcessFunc
		if failFast {
			processFunc = cmd.ProcessSerially
		} else {
			processFunc = cmd.ProcessSeriallyBatchErrors
		}

		if err := processFunc(func(buildSpecWithDeps params.ProductBuildSpecWithDeps, stdout io.Writer) error {
			return Run(buildSpecWithDeps, publisher, almanacInfo, stdout)
		})(buildSpecWithDeps, stdout); err != nil {
			// if publish failed with bulk errors, print nice error message
			if specErrors, ok := err.(*cmd.SpecErrors); ok {
				var parts []string
				for _, v := range specErrors.Errors {
					parts = append(parts, fmt.Sprintf("%v", v))
				}
				return fmt.Errorf(strings.Join(parts, "\n"))
			}
			return err
		}

		return nil
	}, cfg, products, wd, stdout)
}
