####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/cat-sleep-example-${TARGETARCH} /bin/cat-sleep-example
RUN chmod +x /bin/cat-sleep-example

####################################################################################################
# cat-sleep
####################################################################################################
FROM scratch AS cat-sleep
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/cat-sleep-example /bin/cat-sleep-example
ENTRYPOINT [ "/bin/cat-sleep-example" ]
