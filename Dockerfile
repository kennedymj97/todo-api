# Pull image from digest to ensure golang:alpine is not being intercepted.
# Update that hash to that latest by doing docker pull golang:alpine.
FROM golang@sha256:06ba1dae97f2bf560831497f8d459c68ab75cc67bf6fc95d9bd468ac259c9924 as builder
# Install git, add ssl ca certificate and zoneinfo for timezones
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
# Create user so the container is not run by root
RUN adduser -D -g '' appuser
# Could I pull the code from git here? so ADD ./project/ / instead???
ADD . /
WORKDIR /
# Ensure all go dependencies are downloaded
RUN go mod download
# Build the go app statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
# Copy items needed from first stage
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /main /app/
WORKDIR /app
USER appuser
EXPOSE 8080
CMD ["./main"]