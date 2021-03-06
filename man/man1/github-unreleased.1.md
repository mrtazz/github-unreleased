---
title: github-unreleased(1) github-unreleased User Manuals | github-unreleased User Manuals
author: Daniel Schauenberg <d@unwiredcouch.com>
date: REPLACE_DATE
---

# NAME
**github-unreleased** -- find unreleased commits of a GitHub repository

# SYNOPSIS

**github-unreleased** \[options\] \[repo\]

# DESCRIPTION
`github-unreleased` is a tool that will tell you how many commits have been
committed to a repository after its most recent release. It can be passed a
specific repository slug (like mrtazz/github-unreleased) or by default will
check all repositories for the currently authenticated user that have at least
one tag and are not forks.

# OPTIONS

**-h**, **--help**
:    Show a friendly help message.

**--version**
:    Show version.

**--debug**
:    Enable debug mode

**--include-forks**
:    Also show unreleased commits in forks

**--inlcude-repos-without-tags**
:    Also show repositories with no tags

# CONFIGURATION
**github-unreleased** reads an ini configuration file at
~/.github-unreleased.ini and reads configuration values from its **default**
section. The following configuration fields are supported:

**token**
:    the GitHub OAuth token to use


# EXIT STATUS
github-unreleased exits with the following status:

```
 0:   no commits were added after the most recent release
>1:   number of commits that were found after the most recent release
```

In addition to github-unreleased having found 1 unreleased commit, exit status
1 is also used to denote an error in execution happening.

# BUGS
Please file bugs against the issue tracker:
https://github.com/mrtazz/github-unreleased/issues

# SEE ALSO
git(1), hub(1)
