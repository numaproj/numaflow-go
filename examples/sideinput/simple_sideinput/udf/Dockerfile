####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/udf-sideinput-example-${TARGETARCH} /bin/udf-sideinput-example
RUN chmod +x /bin/udf-sideinput-example

####################################################################################################
# udf-sideinput
####################################################################################################
FROM scratch AS udf-sideinput
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/udf-sideinput-example /bin/udf-sideinput-example
ENTRYPOINT [ "/bin/udf-sideinput-example" ]
