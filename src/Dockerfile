ARG GOVERSION=1.26.3

FROM docker.io/library/golang:${GOVERSION}-alpine AS build
ARG GOVERSION
WORKDIR /
COPY ./ /

RUN go mod edit -C main -go $GOVERSION
RUN go mod tidy -C main

RUN go build -C main -o github-exporter main.go

FROM docker.io/library/golang:${GOVERSION}-alpine
ARG GOVERSION

COPY --from=build ./main/github-exporter /github-exporter

WORKDIR /
CMD [ "/github-exporter" ]

