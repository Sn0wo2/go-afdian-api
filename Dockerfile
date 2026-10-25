FROM --platform=$BUILDPLATFORM alpine AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/go-afdian-api /app/go-afdian-api
WORKDIR /app
ENTRYPOINT ["/app/go-afdian-api"]
