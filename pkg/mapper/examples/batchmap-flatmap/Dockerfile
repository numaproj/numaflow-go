####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/flatmap-example /bin/flatmap-example
RUN chmod +x /bin/flatmap-example

####################################################################################################
# flatmap
####################################################################################################
FROM scratch as flatmap
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/flatmap-example /bin/flatmap-example
ENTRYPOINT [ "/bin/flatmap-example" ]
