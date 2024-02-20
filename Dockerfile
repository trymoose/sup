FROM golang:1.22 as builder

WORKDIR /sup
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /sup

FROM scratch

WORKDIR /
COPY --from=builder /sup .
EXPOSE 8080
EXPOSE 9090

ENV DATA_DIR=/data
ENV UPLOAD_ADDRESS=0.0.0.0:9090
ENV UI_ADDRESS=0.0.0.0:8080
ENV SHUTDOWN_TIMEOUT=5s
ENV UPLOAD_TIMEOUT=15s
ENV LIMIT_UPLOAD=0
ENV DESTINATIONS_FILE=/destinations.yml

CMD ["/sup"]