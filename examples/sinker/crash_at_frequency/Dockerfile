####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/crash-at-frequency-${TARGETARCH} /bin/crash-at-frequency
RUN chmod +x /bin/crash-at-frequency

####################################################################################################
# crash-at-frequency
####################################################################################################
FROM scratch AS crash-at-frequency
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/crash-at-frequency /bin/crash-at-frequency
ENTRYPOINT [ "/bin/crash-at-frequency" ]
