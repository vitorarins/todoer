ARG GOVERSION

FROM golang:${GOVERSION} as base

ARG VERSION

WORKDIR /build

COPY . .

RUN go build -o todoer -ldflags "-X main.VersionString=${VERSION}" ./cmd/todoer/todoer.go

# Use two stages only to avoid source code on final image
FROM golang:${GOVERSION}

COPY --from=base /build/todoer /app/todoer

ENTRYPOINT ["/app/todoer"]