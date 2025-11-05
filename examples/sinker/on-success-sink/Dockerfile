####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/on-success-sink-${TARGETARCH} /bin/on-success-sink
RUN chmod +x /bin/on-success-sink

####################################################################################################
# on_success
####################################################################################################
FROM scratch AS on_success
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/on-success-sink /bin/on-success-sink
ENTRYPOINT [ "/bin/on-success-sink" ]
