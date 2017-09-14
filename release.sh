#!/bin/bash

version=$(cat version.txt)
sed -i '' "s/VERSION/${version}/" version/version.go
git commit -am "version bump"
git tag "v${version}"
export GITHUB_TOKEN=$(cat .ghtoken)
goreleaser --rm-dist
sed -i '' "s/${version}/VERSION/" version/version.go
git add .
git commit -am "Built and released ${version}. Updated dist/"
git push