#!/bin/bash

function evil_git_dirty {
  [[ $(git diff --shortstat 2> /dev/null | tail -n1) != "" ]] && echo "*"
}

if [ evil_git_dirty == "*" ]; then
  echo "git index dirty. aborting."
  exit 1
fi

version=$(cat version.txt)
sed -i '' "s/VERSION/${version}/" version/version.go
git add .
git commit --amend --no-edit
git tag "v${version}"
export GITHUB_TOKEN=$(cat .ghtoken)
goreleaser --rm-dist
sed -i '' "s/${version}/VERSION/" version/version.go
git add .
git commit --amend --no-edit
git push