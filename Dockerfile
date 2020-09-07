FROM golang:1.15 as build

LABEL maintainer="Alex <32b3@protonmail.com>"

WORKDIR /goapp

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o app

FROM gcr.io/distroless/base-debian10
COPY --from=build /goapp /

ENV PORT=8000

CMD ["/app"]
