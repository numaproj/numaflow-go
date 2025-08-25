####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/log-example-${TARGETARCH} /bin/log-example
RUN chmod +x /bin/log-example

####################################################################################################
# log
####################################################################################################
FROM scratch AS log
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/log-example /bin/log-example
ENTRYPOINT [ "/bin/log-example" ]
