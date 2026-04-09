FROM docker.io/library/golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -o aimux


FROM docker.io/library/debian:bookworm

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates passwd && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/

COPY --from=builder /app/aimux /app/aimux
COPY --from=builder /app/conf/app.yml /app/conf/app.yml

RUN useradd -m -u 10001 work && chown -R work:work /app
USER work

EXPOSE 8080 8081

ENV APPListenAPI="0.0.0.0:8080" \
    APPListenAdmin="0.0.0.0:8081"

ENTRYPOINT ["/app/aimux"]