FROM golang:1.20.0 AS build

WORKDIR /workdir

COPY go.mod ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o /app ./server/app.go

FROM gcr.io/distroless/base-debian12 AS release

WORKDIR /

COPY --from=build /app /app

EXPOSE 13337

USER nonroot:nonroot

CMD ["/app"]