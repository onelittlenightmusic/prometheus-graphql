FROM golang:1.17.2 AS build

WORKDIR /app

COPY app ./app

COPY go.mod ./
COPY go.sum ./
COPY gqlgen.yml ./
COPY config ./config
COPY graph ./graph

RUN go mod download

RUN pwd && ls && mkdir bin && go build -o bin/server app/server.go 

EXPOSE 8080

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/bin/server bin/server

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["bin/server"]