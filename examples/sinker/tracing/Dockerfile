####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/sink-tracing-example-${TARGETARCH} /bin/sink-tracing-example
RUN chmod +x /bin/sink-tracing-example

####################################################################################################
# sink-tracing-example
####################################################################################################
FROM scratch AS sink-tracing-example
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/sink-tracing-example /bin/sink-tracing-example
ENTRYPOINT [ "/bin/sink-tracing-example" ]
