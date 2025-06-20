####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/time-extraction-filter-example-${TARGETARCH} /bin/time-extraction-filter-example
RUN chmod +x /bin/time-extraction-filter-example

####################################################################################################
# time-extraction-filter
####################################################################################################
FROM scratch AS time-extraction-filter
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/time-extraction-filter-example /bin/time-extraction-filter-example
ENTRYPOINT [ "/bin/time-extraction-filter-example" ]