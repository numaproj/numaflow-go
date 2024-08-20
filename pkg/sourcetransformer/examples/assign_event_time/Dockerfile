####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/assign-event-time-example-${TARGETARCH} /bin/assign-event-time-example
RUN chmod +x /bin/assign-event-time-example

####################################################################################################
# assign-event-time
####################################################################################################
FROM scratch AS assign-event-time
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/assign-event-time-example /bin/assign-event-time-example
ENTRYPOINT [ "/bin/assign-event-time-example" ]
