####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/serving-inmem-store-${TARGETARCH} /bin/serving-inmem-store
RUN chmod +x /bin/serving-inmem-store

####################################################################################################
# flatmap
####################################################################################################
FROM scratch AS memory_store
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/serving-inmem-store /bin/serving-inmem-store
ENTRYPOINT [ "/bin/serving-inmem-store" ]
