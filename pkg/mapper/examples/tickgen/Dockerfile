####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/tickgen-example /bin/tickgen-example
RUN chmod +x /bin/tickgen-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch as tickgen
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/tickgen-example /bin/tickgen-example
ENTRYPOINT [ "/bin/tickgen-example" ]
