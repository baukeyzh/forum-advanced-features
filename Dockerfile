FROM golang:1.18.1-alpine AS build
RUN apk --no-cache add ca-certificates git
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY  . .
RUN apk add build-base
RUN go build -o forum cmd/main.go
FROM alpine:latest
WORKDIR /
COPY --from=build /app .
EXPOSE 8082
CMD ["/forum"]