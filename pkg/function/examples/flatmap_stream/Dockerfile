####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/flatmap-stream-example /bin/flatmap-stream-example
RUN chmod +x /bin/flatmap-stream-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch as flatmap_stream
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/flatmap-stream-example /bin/flatmap-stream-example
ENTRYPOINT [ "/bin/flatmap-stream-example" ]
