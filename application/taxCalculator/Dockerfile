FROM golang:1.24.4-alpine AS build
LABEL authors="Hesam Eskandari"

RUN apk update && apk upgrade
WORKDIR /usr/src
RUN mkdir -p application
COPY . .
WORKDIR /usr/src/application/taxCalculator
RUN go mod tidy
RUN go build main.go


FROM golang:1.24.4-alpine AS server
RUN apk update && apk upgrade
RUN addgroup -S leafgp && adduser -S leaf -G leafgp

WORKDIR /app
RUN mkdir -p .logs && chown leaf:leafgp -R .logs
USER leaf

COPY --from=build --chown=leaf:leafgp --chmod=100 /usr/src/application/taxCalculator/main ./main
COPY --from=build --chown=leaf:leafgp --chmod=400 /usr/src/application/taxCalculator/service/fileStorage/data ./application/taxCalculator/service/fileStorage/data

ENTRYPOINT ["./main"]
