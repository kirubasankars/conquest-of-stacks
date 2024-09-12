#!/bin/sh

set -e

export GIN_MODE=release
export CENTRIFUGO_API_KEY=$(jq -r '.api_key' /tmp/centrifugo.json)
export HMAC_SECRET_KEY=$(jq -r '.token_hmac_secret_key' /tmp/centrifugo.json)

/cos/cos