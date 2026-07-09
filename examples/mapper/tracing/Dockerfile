####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/map-tracing-example-${TARGETARCH} /bin/map-tracing-example
RUN chmod +x /bin/map-tracing-example

####################################################################################################
# map-tracing-example
####################################################################################################
FROM scratch AS map-tracing-example
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/map-tracing-example /bin/map-tracing-example
ENTRYPOINT [ "/bin/map-tracing-example" ]
