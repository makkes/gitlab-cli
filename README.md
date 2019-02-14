# GitLab Command Line Interface

[![Build Status](https://travis-ci.org/makkes/gitlab-cli.svg?branch=master)](https://travis-ci.org/makkes/gitlab-cli)

The GitLab Command Line Interface (CLI) is a cross-platform command line utility
that provides a user-friendly yet powerful way to query information from your
GitLab repos.

![](./demo.gif "GitLab CLI Demo Video")

## Installation and Usage

```
go get github.com/makkes/gitlab-cli
```

or grab the binary of the [most current
release](https://github.com/makkes/gitlab-cli/releases).

All commands of gitlab-cli currently require that you are authenticated. To do
so you issue `gitlab-cli login YOUR_TOKEN`. You obtain a personal access token
at https://gitlab.com/profile/personal_access_tokens.

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
. <(gitlab-cli completion)
```

To have completion set up for you automatically just copy and paste the line
from above into your `~/.bashrc` or `~/.profile`.

## License

This software is distributed under the BSD 2-Clause License, see
[LICENSE](LICENSE) for more information.

