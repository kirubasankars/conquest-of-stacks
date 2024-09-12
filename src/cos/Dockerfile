FROM fedora as builder

RUN dnf -y install go jq ca-certificates && \
    mkdir -p cos

COPY go.mod .
COPY go.sum .
RUN GOOS=linux GOARCH=arm64 go mod download

COPY . cos
WORKDIR cos
RUN GOOS=linux GOARCH=arm64 go build -o cos
RUN chmod +x /cos/entrypoint.sh

FROM fedora

RUN dnf -y install jq ca-certificates && mkdir -p /cos

COPY --from=builder /cos/entrypoint.sh /cos/entrypoint.sh
COPY --from=builder /cos/cos /cos/cos
COPY --from=builder /cos/www /cos/www

WORKDIR /cos

ENTRYPOINT ["/cos/entrypoint.sh"]