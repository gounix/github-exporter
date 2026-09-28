ARG GOVERSION
ARG IMAGE_VERSION

FROM docker.io/library/golang:${GOVERSION}-alpine AS build
ARG GOVERSION
ARG IMAGE_VERSION
WORKDIR /
COPY ./ /

RUN go mod edit -C main -go $GOVERSION
RUN go mod tidy -C main
RUN sed -i s/Development-version/${IMAGE_VERSION}/ /main/main.go
RUN sed -i s/Golang-version/${GOVERSION}/ /main/main.go


RUN go build -C main -o github-exporter main.go

FROM docker.io/library/golang:${GOVERSION}-alpine
ARG GOVERSION
ARG IMAGE_VERSION

COPY --from=build ./main/github-exporter /github-exporter

WORKDIR /
CMD [ "/github-exporter" ]

