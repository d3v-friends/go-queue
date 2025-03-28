#!/bin/zsh

BASE_TAG=$1
TAG="v$BASE_TAG"
ORIGIN="origin"

echo "tag=$TAG, base_tag=$BASE_TAG, origin=$ORIGIN"

git tag -d "$TAG"
git push -d "$ORIGIN" "$TAG"
git tag "$TAG"
git push --tags "$ORIGIN" "$TAG"

echo "upload tag done!"
