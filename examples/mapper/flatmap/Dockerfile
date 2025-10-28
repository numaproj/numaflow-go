####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/flatmap-example-${TARGETARCH} /bin/flatmap-example
RUN chmod +x /bin/flatmap-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch AS flatmap
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/flatmap-example /bin/flatmap-example
ENTRYPOINT [ "/bin/flatmap-example" ]
