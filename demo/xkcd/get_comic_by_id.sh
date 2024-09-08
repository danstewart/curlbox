#!/usr/bin/env bash

comic_id="$1"
if [[ -z $comic_id ]]; then
    echo "Usage: $0 <comic_id>"
    exit 1
fi

curl ${PROTOCOL}://xkcd.com/${comic_id}/info.0.json

