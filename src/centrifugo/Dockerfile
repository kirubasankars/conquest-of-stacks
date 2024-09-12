FROM centrifugo/centrifugo

COPY entrypoint.sh /entrypoint.sh
COPY additional_config.json /tmp/additional_config.json

RUN apk --no-cache add jq && chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]