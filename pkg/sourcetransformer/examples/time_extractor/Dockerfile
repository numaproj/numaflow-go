####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/time-extractor-example-${TARGETARCH} /bin/time-extractor-example
RUN chmod +x /bin/time-extractor-example

####################################################################################################
# time-extractor
####################################################################################################
FROM scratch AS time-extractor
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/time-extractor-example /bin/time-extractor-example
ENTRYPOINT [ "/bin/time-extractor-example" ]