####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/even-odd-key-example /bin/even-odd-key-example
RUN chmod +x /bin/even-odd-key-example

####################################################################################################
# even-odd-key
####################################################################################################
FROM scratch as even-odd-key
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/even-odd-key-example /bin/even-odd-key-example
ENTRYPOINT [ "/bin/even-odd-key-example" ]
