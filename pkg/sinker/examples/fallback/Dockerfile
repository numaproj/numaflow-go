####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/fallback-example /bin/fallback-example
RUN chmod +x /bin/fallback-example

####################################################################################################
# fallback
####################################################################################################
FROM scratch as fallback
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/fallback-example /bin/fallback-example
ENTRYPOINT [ "/bin/fallback-example" ]
