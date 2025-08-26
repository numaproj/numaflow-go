####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/redis-sink-${TARGETARCH} /bin/redis-sink
RUN chmod +x /bin/redis-sink

####################################################################################################
# redis
####################################################################################################
FROM scratch AS redis
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/redis-sink /bin/redis-sink
ENTRYPOINT [ "/bin/redis-sink" ]
