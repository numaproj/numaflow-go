####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/event-time-extractor-example-${TARGETARCH} /bin/event-time-extractor-example
RUN chmod +x /bin/event-time-extractor-example

####################################################################################################
# event-time-extractor
####################################################################################################
FROM scratch AS event-time-extractor
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/event-time-extractor-example /bin/event-time-extractor-example
ENTRYPOINT [ "/bin/event-time-extractor-example" ]