####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/on-success-log-${TARGETARCH} /bin/on-success-log
RUN chmod +x /bin/on-success-log

####################################################################################################
# on_success_log
####################################################################################################
FROM scratch AS on_success_log
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/on-success-log /bin/on-success-log
ENTRYPOINT [ "/bin/on-success-log" ]
