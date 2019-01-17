# git subcommand for semver 

### Usage
~~~
go get github.com/liangzuobin/git-semver
~~~

Current semver
~~~
➜  git-semver master ✗ git semver current
current semver v2.0.0
~~~

Patch
~~~
➜  git-semver master ✗ git semver patch
current version: v2.0.1
~~~

Minor
~~~
➜  git-semver master ✗ git semver minor
current version: v2.1.0
~~~

Major
~~~
➜  git-semver master ✗ git semver major
current version: v3.0.0
~~~

BTW, rm all your local tags
~~~
git tag -d $(git tag -l)
~~~