####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/simple-example /bin/simple-example
RUN chmod +x /bin/simple-example

####################################################################################################
# simple-source
####################################################################################################
FROM scratch as simple-source
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/simple-source-example /bin/simple-source-example
ENTRYPOINT [ "/bin/simple-source-example" ]
