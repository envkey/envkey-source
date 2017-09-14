#!/bin/bash

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