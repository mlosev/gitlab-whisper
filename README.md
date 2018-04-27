# gitlab-whisper

Have you ever tried to clone or update multiple projects at once from private Gitlab instance?\
I decided to do that once, and, as a result, this project was born.

One might say this task can be easily solved with a simple bash script and would be right.\
But there would be no fun in writing boring bash scripts when we have such a beautiful language as Go
and maybe some spare time.

To tell the long story short, I created this tool with the primary goal of helping me
with cloning (and subsequently updating) a bunch of projects from our private Gitlab instance to my local machine.\
Maybe it might be useful for someone else too, so I'd better share it just in case.

Build
-----

How to build this tool? The answer is simple - you will need "dep" tool to get all dependencies and, of course,
 working go compiler.

There is a simple Makefile, so the whole sequence of commands after cloning git repository is the following:
- `make vendor`
- `make build`

Also, I have prepared some binaries on [releases page](https://github.com/mlosev/tmux-ssh/releases),
so feel free to grab it from there.

Run
---

Here we assume that you have `git` binary installed.

To run you will need two environment variables exported:
- `GITLAB_API_ENDPOINT`
- `GITLAB_API_PRIVATE_TOKEN`

where `GITLAB_API_ENDPOINT` contains URL of Gitlab API endpoint (e.g., https://gitlab.com/api/v4)
and `GITLAB_API_PRIVATE_TOKEN` is a personal token, issued via the web interface of Gitlab instance.

The command of prime interest is `projects sync` and, by default, it works in "dry run" mode
(no local filesystem is touched).\
In this mode, it prints a list of projects, which are going to be cloned or updated
(if the project folder already exists).\
To perform synchronization you need to pass `-r` option and wait for `git` to clone or fetch all the projects,
member of which you are.

By default, projects are cloned to a path defined in GOPATH environment variable, but this could be easily overridden
via `-d` option.

To avoid collisions of projects names directory structure repeats that of `GOPATH`.\
For example:

```
GOPATH/
   |-src
      |-github.com
         |-mlosev
            |-gitlab-whisper
```

This structure has been proven as quite convenient, and in fact, this is entirely a monorepo.

