####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/flatmap-stream-example-${TARGETARCH} /bin/flatmap-stream-example
RUN chmod +x /bin/flatmap-stream-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch AS flatmap_stream
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/flatmap-stream-example /bin/flatmap-stream-example
ENTRYPOINT [ "/bin/flatmap-stream-example" ]
