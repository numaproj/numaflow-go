####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/reduce-sideinput-example /bin/reduce-sideinput-example
RUN chmod +x /bin/reduce-sideinput-example

####################################################################################################
# sideinput
####################################################################################################
FROM scratch as sideinput
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/reduce-sideinput-example /bin/reduce-sideinput-example
ENTRYPOINT [ "/bin/reduce-sideinput-example" ]
