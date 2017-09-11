#!/bin/bash

git tag v$(cat version.txt)
export GITHUB_TOKEN=$(cat .ghtoken)
goreleaser --rm-dist