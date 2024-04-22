####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/e2e-map-sideinput-example /bin/e2e-map-sideinput-example
RUN chmod +x /bin/e2e-map-sideinput-example

####################################################################################################
# sideinput
####################################################################################################
FROM scratch as sideinput-example
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/e2e-map-sideinput-example /bin/e2e-map-sideinput-example
ENTRYPOINT [ "/bin/e2e-map-sideinput-example" ]
