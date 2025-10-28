####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/serving-redis-store-${TARGETARCH} /bin/serving-redis-store
RUN chmod +x /bin/serving-redis-store

####################################################################################################
# flatmap
####################################################################################################
FROM scratch AS redis-store
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/serving-redis-store /bin/serving-redis-store
ENTRYPOINT [ "/bin/serving-redis-store" ]
