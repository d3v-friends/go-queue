TAG=$1;

git tag -d "$TAG"
git push -d origin "$TAG"
git tag "$TAG"
git push origin "$TAG"
