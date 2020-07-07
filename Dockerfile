FROM golang:1.13-alpine3.11 AS build
RUN echo -e "http://nl.alpinelinux.org/alpine/v3.5/main\nhttp://nl.alpinelinux.org/alpine/v3.5/community" > /etc/apk/repositories
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/social-network
COPY . .

COPY ./web /go/bin/
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app cmd/main.go

FROM alpine:3.11
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 3000
CMD ["app"]
