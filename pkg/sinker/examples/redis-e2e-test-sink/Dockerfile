####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/redis-e2e-test-sink /bin/redis-e2e-test-sink
RUN chmod +x /bin/redis-e2e-test-sink

####################################################################################################
# redis
####################################################################################################
FROM scratch as redis
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/redis-e2e-test-sink /bin/redis-e2e-test-sink
ENTRYPOINT [ "/bin/redis-e2e-test-sink" ]
