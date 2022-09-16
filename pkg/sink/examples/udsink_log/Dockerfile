####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/go-udsink-example /bin/go-udsink-example
RUN chmod +x /bin/go-udsink-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch as flatmap
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/go-udsink-example /bin/go-udsink-example
ENTRYPOINT [ "/bin/go-udsink-example" ]
