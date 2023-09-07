REPO=$1
PROJECT=$2
BUILDER_TAG="$REPO/$PROJECT:latest"

docker build \
  --platform linux/amd64 \
  -t "$BUILDER_TAG" \
  -f ./docker/rabbitMq/dockerfile .

docker push "$BUILDER_TAG"
