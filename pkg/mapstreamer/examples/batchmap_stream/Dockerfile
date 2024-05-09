####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/batchmap-stream-example /bin/batchmap-stream-example
RUN chmod +x /bin/batchmap-stream-example

####################################################################################################
# batchmap
####################################################################################################
FROM scratch as batchmap_stream
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/batchmap-stream-example /bin/batchmap-stream-example
ENTRYPOINT [ "/bin/batchmap-stream-example" ]
