####################################################################################################
# base
####################################################################################################
FROM alpine:3.17.1 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/event-time-filter-example /bin/event-time-filter-example
RUN chmod +x /bin/event-time-filter-example

####################################################################################################
# event-time-filter
####################################################################################################
FROM scratch as event-time-filter
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/event-time-filter-example /bin/event-time-filter-example
ENTRYPOINT [ "/bin/event-time-filter-example" ]
