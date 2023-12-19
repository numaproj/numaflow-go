####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/counter-example /bin/counter-example
RUN chmod +x /bin/counter-example

####################################################################################################
# counter
####################################################################################################
FROM scratch as counter
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/counter-example /bin/counter-example
ENTRYPOINT [ "/bin/counter-example" ]
