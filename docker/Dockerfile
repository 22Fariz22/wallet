FROM golang:1.23-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download

FROM golang:1.23-alpine as runner

COPY --from=builder ./app ./app

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app
ENV config=docker

EXPOSE 8080

ENTRYPOINT CompileDaemon --build="go build cmd/api/main.go" --command=./main




