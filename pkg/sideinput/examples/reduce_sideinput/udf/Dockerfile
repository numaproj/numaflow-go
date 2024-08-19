####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/reduce-sideinput-udf-${TARGETARCH} /bin/reduce-sideinput-udf
RUN chmod +x /bin/reduce-sideinput-udf

####################################################################################################
# udf-sideinput
####################################################################################################
FROM scratch AS sideinput-udf
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/reduce-sideinput-udf /bin/reduce-sideinput-udf
ENTRYPOINT [ "/bin/reduce-sideinput-udf" ]