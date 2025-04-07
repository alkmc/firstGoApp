FROM golang:1.24 as builder

LABEL maintainer="Alex <github.com/alkmc>"

WORKDIR /goapp

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/base-debian11
COPY --from=builder /goapp /

ENV PORT=8000

CMD ["/app"]
