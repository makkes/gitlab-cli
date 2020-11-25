# GitLab Command Line Interface

![Build Status](https://github.com/makkes/gitlab-cli/workflows/Test/badge.svg)
[![Build Status](https://travis-ci.org/makkes/gitlab-cli.svg?branch=master)](https://travis-ci.org/makkes/gitlab-cli)

The GitLab Command Line Interface (CLI) is a cross-platform command line utility
that provides a user-friendly yet powerful way to query information from your
GitLab repos.

## Installation and Usage

The easiest installation method is to use the installation script:

```
# install latest version into /usr/local/bin/
curl -sSfL https://raw.githubusercontent.com/makkes/gitlab-cli/master/install.sh | sh -s 

# install latest version into ~/bin/
curl -sSfL https://raw.githubusercontent.com/makkes/gitlab-cli/master/install.sh | sh -s -- -b ~/bin

# install v3.6.3 into /usr/local/bin
curl -sSfL https://raw.githubusercontent.com/makkes/gitlab-cli/master/install.sh | sh -s v3.6.3
```

If that script doesn't work for you and you have a Go environment set up you can
use this command:

```
go get github.com/makkes/gitlab-cli
```

As a last resort just manually grab the binary of the [most current
release](https://github.com/makkes/gitlab-cli/releases).

All commands of gitlab-cli currently require that you are authenticated. To do
so you issue `gitlab login YOUR_TOKEN`. You obtain a personal access token
at https://gitlab.com/profile/personal_access_tokens. To make use of all of
gitlab-cli's features you need to grant api, read_user, read_repository and
read_registry scopes.

## Updating the CLI

Since version 3.6 the CLI has an `update` command that you can use to update the
CLI's version so you don't have to download the latest release every time.

Running
```
gitlab update
```
will update your version to the latest stable release of the current **major** release, so it would update from e.g. 3.7.1 to 3.8.0 but not to 4.0.0. If you'd like to upgrade to the next major release, provide the `--major` flag (available since 3.7.2):
```
gitlab update --major
```

### Dry-run and pre-release updates

If you would like to just check for availability of a new version, use the `--dry-run` flag:
```
gitlab update --dry-run
```

For the brave companions, there's also the `--pre` flag (available since 3.7.3) which will update to the next pre-release (i.e. alpha, beta or release candidate) version:
```
gitlab update --pre
```

## Commands

Currently GitLab CLI supports these commands:

* `projects`: List all your projects
* `project`:  List details about a project by ID or name
* `project create`: Create a new project
* `var`: Manage project variables
* `pipelines`: List pipelines of a project
* `pipeline`: List details of a pipeline
* `issues`: List all issues of a project
* `issue`: Manage issues
* `status`: Display the current configuration of GitLab CLI

## Bash Completion

You can get your Bash to complete GitLab CLI commands very easily: Just type the
following line in your shell:

```sh
. <(gitlab completion)
```

To have completion set up for you automatically just copy and paste the line
from above into your `~/.bashrc` or `~/.profile`.

## License

This software is distributed under the BSD 2-Clause License, see
[LICENSE](LICENSE) for more information.

