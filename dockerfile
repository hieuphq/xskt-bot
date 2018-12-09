FROM golang:latest 

WORKDIR /go/src/github.com/hieuphq/xskt-bot 
ADD . /go/src/github.com/hieuphq/xskt-bot
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api cmd/server/main.go

FROM alpine
WORKDIR /
COPY --from=0 /go/src/github.com/hieuphq/xskt-bot/api api

CMD ["api"]