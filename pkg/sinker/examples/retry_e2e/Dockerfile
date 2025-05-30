####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/retry-e2e-${TARGETARCH} /bin/retry-e2e
RUN chmod +x /bin/retry-e2e

####################################################################################################
# retry-e2e
####################################################################################################
FROM scratch AS retry-e2e
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/retry-e2e /bin/retry-e2e
ENTRYPOINT [ "/bin/retry-e2e" ]