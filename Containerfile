ARG GO_VERSION=1.24-alpine
FROM docker.io/golang:${GO_VERSION} as builder

# Install git + SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates
# Create user worker
RUN adduser -D -g '' worker
COPY . /build
WORKDIR /build

ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

#build the binary
RUN \
    go build \
    -ldflags "-d -s -w -extldflags \"-static\"" \
    -a -tags netgo -installsuffix netgo \
    -o /go/bin/terraform-sops-backend

# STEP 2 package the result image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/terraform-sops-backend /bin/terraform-sops-backend

USER worker
# WORKDIR /data
ENTRYPOINT ["/bin/terraform-sops-backend", "start"]
