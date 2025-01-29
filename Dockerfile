# https://hackmd.io/@_Cl3aSMeQf2V5wwqxEyfwg/SkR0yzU5s
FROM golang:1.23-alpine
WORKDIR /go/src/app

COPY . .
# ENV CGO_ENABLED=1

RUN go build -o main .

CMD ["./main"]
