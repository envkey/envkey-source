#!/bin/bash

version=$(cat version.txt)
sed -i "s/VERSION/${version}/g" version/version.go
git tag "v${version}"
export GITHUB_TOKEN=$(cat .ghtoken)
goreleaser --rm-dist
git add .
git commit -am "Built and released ${version}. Updated dist/"
git push