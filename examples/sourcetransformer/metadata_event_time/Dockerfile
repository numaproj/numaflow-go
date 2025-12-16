####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/metadata-event-time-example-${TARGETARCH} /bin/metadata-event-time-example
RUN chmod +x /bin/metadata-event-time-example

####################################################################################################
# metadata-event-time
####################################################################################################
FROM scratch AS metadata-event-time
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/metadata-event-time-example /bin/metadata-event-time-example
ENTRYPOINT [ "/bin/metadata-event-time-example" ]
