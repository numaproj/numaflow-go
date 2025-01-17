####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/slow-cat-example-${TARGETARCH} /bin/slow-cat-example
RUN chmod +x /bin/slow-cat-example

####################################################################################################
# slow-cat
####################################################################################################
FROM scratch AS slow-cat
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/slow-cat-example /bin/slow-cat-example
ENTRYPOINT [ "/bin/slow-cat-example" ]
