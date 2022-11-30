FROM node:18.12.1-alpine as node

FROM golang:1.19.3-bullseye

# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get --no-install-recommends -y install gnupg2 ca-certificates sqlite3 libsqlite3-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists

COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/include/node /usr/local/include/node
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
RUN ln -s /usr/local/bin/node /usr/local/bin/nodejs && \
    ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
    ln -s /opt/yarn/bin/yarn /usr/local/bin/yarn

WORKDIR /go-stations

RUN go install github.com/cweill/gotests/gotests@latest && \
    go install github.com/fatih/gomodifytags@latest && \
    go install github.com/josharian/impl@latest && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install golang.org/x/tools/gopls@latest

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /bin/go-stations
