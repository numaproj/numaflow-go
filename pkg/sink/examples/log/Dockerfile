####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/log-example /bin/log-example
RUN chmod +x /bin/log-example

####################################################################################################
# log
####################################################################################################
FROM scratch as log
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/log-example /bin/log-example
ENTRYPOINT [ "/bin/log-example" ]
