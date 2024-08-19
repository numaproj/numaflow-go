####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/map-sideinput-udf-${TARGETARCH} /bin/map-sideinput-udf
RUN chmod +x /bin/map-sideinput-udf

####################################################################################################
# udf-sideinput
####################################################################################################
FROM scratch AS sideinput-udf
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/map-sideinput-udf /bin/map-sideinput-udf
ENTRYPOINT [ "/bin/map-sideinput-udf" ]