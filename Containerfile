FROM golang:1.24 AS builder
WORKDIR /build

# Cache module downloads.
COPY go.mod /build/
RUN go mod download

# Copy source and build.
COPY . /build
ENV CGO_ENABLED=0
RUN go build -v -o /build/bin/test-check ./cmd/test-check

# Create a non-root user.
RUN groupadd -g 999 user && \
    useradd -r -u 999 -g user user

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bin/test-check /app/test-check
USER user
ENTRYPOINT ["/app/test-check"]
