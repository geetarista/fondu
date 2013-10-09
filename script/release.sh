#!/usr/bin/env bash

set -e

VERSION=$1

DATA="{\"tag_name\":\"${VERSION}\",\"target_commitish\":\"master\",\"name\":\"v${VERSION}\"}"
DATA="'${DATA}'"
HEADERS="'Accept: application/vnd.github.manifold-preview'"
TOKEN=$(cat ~/.github_token)
URL="https://api.github.com/repos/geetarista/fondu/releases"
USER="'geetarista:${TOKEN}'"

CMD="curl -s -X POST -u $USER -H $HEADERS -d $DATA $URL"
JSON=$(eval $CMD)
ID=$(echo $JSON | python -c 'import sys, json; print json.load(sys.stdin)[sys.argv[1]]' id)

pushd bin
for file in *; do
  ACCEPT="'Accept: application/vnd.github.manifold-preview'"
  CONTENT_TYPE="'Content-Type: application/octet-stream'"
  DATA="{\"name\":\"${file}\"}"
  DATA="'${DATA}'"
  URL="https://uploads.github.com/repos/geetarista/fondu/releases/${ID}/assets?name=${file}"
  CMD="curl -s -X POST -u $USER -H $ACCEPT -H $CONTENT_TYPE -d $DATA $URL"
  JSON=$(eval $CMD)
done
popd
