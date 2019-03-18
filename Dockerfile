FROM golang:1.12.1-alpine3.9 as builder

RUN apk --no-cache add make git && rm -rf /var/cache/apk/*
ARG VERSION
WORKDIR /go/src/deploy-wizard

ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make VERSION=$VERSION build

FROM alpine:3.9
RUN apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/deploy-wizard/deploy-wizard /bin
CMD ["/bin/deploy-wizard"]
