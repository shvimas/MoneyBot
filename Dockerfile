FROM golang:alpine

WORKDIR /go/src/MoneyBot

COPY . .

RUN apk add --no-cache git

WORKDIR /go/src/MoneyBot/money-bot
RUN go get -v -d
RUN go install -v

RUN apk del git

RUN export PATH="$GOPATH/bin:$PATH"

ENTRYPOINT ["money-bot"]