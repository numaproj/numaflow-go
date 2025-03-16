####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/streamsorter-${TARGETARCH} /bin/streamsorter
RUN chmod +x /bin/streamsorter

####################################################################################################
# accumulator
####################################################################################################
FROM scratch AS streamsorter
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/streamsorter /bin/streamsorter
ENTRYPOINT [ "/bin/streamsorter" ]
