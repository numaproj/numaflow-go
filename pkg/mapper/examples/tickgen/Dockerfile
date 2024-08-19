####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/tickgen-example-${TARGETARCH} /bin/tickgen-example
RUN chmod +x /bin/tickgen-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch AS tickgen
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/tickgen-example /bin/tickgen-example
ENTRYPOINT [ "/bin/tickgen-example" ]
