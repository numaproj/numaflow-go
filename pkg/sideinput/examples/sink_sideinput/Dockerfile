####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/redis-sink-with-sideinput-${TARGETARCH} /bin/redis-sink-with-sideinput
RUN chmod +x /bin/redis-sink-with-sideinput

####################################################################################################
# redis
####################################################################################################
FROM scratch AS sideinput
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/redis-sink-with-sideinput /bin/redis-sink-with-sideinput
ENTRYPOINT [ "/bin/redis-sink-with-sideinput" ]
