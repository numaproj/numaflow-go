####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/metadata-cat-example-${TARGETARCH} /bin/metadata-cat-example
RUN chmod +x /bin/metadata-cat-example

####################################################################################################
# metadata-cat-example
####################################################################################################
FROM scratch AS metadata-cat-example
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/metadata-cat-example /bin/metadata-cat-example
ENTRYPOINT [ "/bin/metadata-cat-example" ]
