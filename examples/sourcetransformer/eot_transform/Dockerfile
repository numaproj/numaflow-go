####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/eot-transform-${TARGETARCH} /bin/eot-transform
RUN chmod +x /bin/eot-transform

####################################################################################################
# eot-transform
####################################################################################################
FROM scratch AS eot-transform
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/eot-transform /bin/eot-transform
ENTRYPOINT [ "/bin/eot-transform" ]
