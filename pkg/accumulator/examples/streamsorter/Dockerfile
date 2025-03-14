####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/accumulator-example-${TARGETARCH} /bin/accumulator-example
RUN chmod +x /bin/accumulator-example

####################################################################################################
# accumulator
####################################################################################################
FROM scratch AS accumulator
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/accumulator-example /bin/accumulator-example
ENTRYPOINT [ "/bin/accumulator-example" ]
