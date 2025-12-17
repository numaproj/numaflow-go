####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/metadata-log-example-${TARGETARCH} /bin/metadata-log-example
RUN chmod +x /bin/metadata-log-example

####################################################################################################
# metadata-log-example
####################################################################################################
FROM scratch AS metadata-log-example
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/metadata-log-example /bin/metadata-log-example
ENTRYPOINT [ "/bin/metadata-log-example" ]
