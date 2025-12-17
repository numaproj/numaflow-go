####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/simple-source-with-metadata-example-${TARGETARCH} /bin/simple-source-with-metadata-example
RUN chmod +x /bin/simple-source-with-metadata-example

####################################################################################################
# simple-source-with-metadata
####################################################################################################
FROM scratch AS simple-source-with-metadata
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/simple-source-with-metadata-example /bin/simple-source-with-metadata-example
ENTRYPOINT [ "/bin/simple-source-with-metadata-example" ]
