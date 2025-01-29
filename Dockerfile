# https://hackmd.io/@_Cl3aSMeQf2V5wwqxEyfwg/SkR0yzU5s
FROM golang:1.23-alpine
WORKDIR /go/src/app

COPY . .
# ENV CGO_ENABLED=1

# Install grype
RUN apk add --no-cache curl
RUN curl -sSfL https://raw.githubusercontent.com/anchore/grype/main/install.sh | sh -s -- -b /usr/local/bin

RUN go build -o main .

CMD ["./main"]
