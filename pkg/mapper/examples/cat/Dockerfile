####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/cat-example-${TARGETARCH} /bin/cat-example
RUN chmod +x /bin/cat-example

####################################################################################################
# cat
####################################################################################################
FROM scratch AS cat
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/cat-example /bin/cat-example
ENTRYPOINT [ "/bin/cat-example" ]
