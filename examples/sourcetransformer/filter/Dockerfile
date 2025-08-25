####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/filter-example-${TARGETARCH} /bin/filter-example
RUN chmod +x /bin/filter-example

####################################################################################################
# filter
####################################################################################################
FROM scratch AS filter
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/filter-example /bin/filter-example
ENTRYPOINT [ "/bin/filter-example" ]
