# github-unreleased

[![Build Status](https://travis-ci.org/mrtazz/github-unreleased.svg?branch=master)](https://travis-ci.org/mrtazz/github-unreleased)
[![Coverage Status](https://coveralls.io/repos/mrtazz/github-unreleased/badge.svg?branch=master&service=github)](https://coveralls.io/github/mrtazz/github-unreleased?branch=master)
[![Packagecloud](https://img.shields.io/badge/packagecloud-available-brightgreen.svg)](https://packagecloud.io/mrtazz/github-unreleased)
[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg)](http://opensource.org/licenses/MIT)

## Overview
A simple command line tool to show you commits of a GitHub repo that have been
added after the most recent release.

## Usage
`github-unreleased` supports a couple of command line arguments. But in
general you either give it a repo to check or by default it will check all
repos that are owned by the currently authenticated user, have at least one
tag and are not forks.

```
% github-unreleased --help
github-unreleased.

  Usage:
  github-unreleased [options] [<repository>]
  github-unreleased -h | --help
  github-unreleased --version

  Options:
  -h --help                       Show this screen.
  --version                       Show version.
  --debug                         Enable debug mode
  --include-forks                 Also show unreleased commits in forks
  --include-repos-without-tags    Also show repositories with no releases

```

```
% github-unreleased mrtazz/snyder

 ==> There are 4 commits since tag "0.4.4" (12/24/2015) for "mrtazz/snyder"

     SHA             AUTHOR                    MESSAGE                                                     URL

  f50dedac3c   Daniel Schauenberg   add doxygen and MIT license      https://github.com/mrtazz/snyder/commit/f50dedac3cbf6dddc02448c04085aa7b1a479de0
                                    badge to README
  55101b0ddb   Keyur                Run ldconfig after package       https://github.com/mrtazz/snyder/commit/55101b0ddb694aa6a10687570fcf033a2ee8753c
                                    install/uninstall
  3fddba5e10   Daniel Schauenberg   Merge pull request #5            https://github.com/mrtazz/snyder/commit/3fddba5e103df1e9977a0d46859f18f0e2287bbf
                                    from keyurdg/master  Run
                                    ldconfig after package
                                    install/uninstall
  a90ddfae95   Daniel Schauenberg   fix docs build for gh-pages      https://github.com/mrtazz/snyder/commit/a90ddfae9500f1c3413b872e35c70abdb47c67ac
                                    Jekyll upgrade
```

## Configuration
`github-unreleased` checks for an ini file in `~/.github-unreleased.ini` which
contains the GitHub OAuth token as follows:

```
[default]
token = 1234foobla
```
