####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/cat-with-metadata-example-${TARGETARCH} /bin/cat-with-metadata-example
RUN chmod +x /bin/cat-with-metadata-example

####################################################################################################
# cat
####################################################################################################
FROM scratch AS cat-with-metadata
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/cat-with-metadata-example /bin/cat-with-metadata-example
ENTRYPOINT [ "/bin/cat-with-metadata-example" ]
