####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/batchmap-flatmap /bin/batchmap-flatmap
RUN chmod +x /bin/batchmap-flatmap

####################################################################################################
# batchmap-flatmap
####################################################################################################
FROM scratch as batch-flatmap
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/batchmap-flatmap /bin/batchmap-flatmap
ENTRYPOINT [ "/bin/batchmap-flatmap" ]
