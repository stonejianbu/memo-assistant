# Build Stage
FROM golang:1.24.3 AS build-stage

ADD . /src/server
WORKDIR /src/server

ARG GIT_COMMIT
ARG APP_NAME

ARG VERSION
ENV VERSION=${VERSION}
ENV GOPROXY="https://goproxy.cn,direct"

COPY go.mod go.sum ./
RUN go mod download
RUN go build ${ENABLE_RACE} -ldflags "-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X main.Version=${VERSION}" -o bin/server main.go

# Final Stage
FROM alpine:latest AS final

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache tcpdump lsof net-tools tzdata curl dumb-init libc6-compat

WORKDIR /app

COPY --from=build-stage /src/server/bin/server /app
COPY --from=build-stage /src/server/config/config.yaml /app

RUN chmod +x /app/server

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/app/server", "-conf", "/app/config.yaml"]