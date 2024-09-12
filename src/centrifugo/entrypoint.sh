#!/bin/sh

set -e
jq -s '.[0] * .[1]' /tmp/config.json /tmp/additional_config.json > /centrifugo/config.json
centrifugo