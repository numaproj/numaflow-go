####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/serve-example-${TARGETARCH} /bin/serve-example
RUN chmod +x /bin/serve-example

####################################################################################################
# serve
####################################################################################################
FROM scratch AS serve
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/serve-example /bin/serve-example
ENTRYPOINT [ "/bin/serve-example" ]
