FROM golang:alpine AS build

WORKDIR /app
COPY go.* .
RUN go mod download

COPY . .
RUN go build -o api cmd/main.go

FROM build AS test

CMD go test ./... -cover -coverprofile=./coverage/coverage.txt -covermode=atomic

FROM alpine AS release

COPY --from=build app/api .
COPY ./data/migrations ./migrations
EXPOSE 8080
ENTRYPOINT ["./api"]