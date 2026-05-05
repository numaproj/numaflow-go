####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/no-response-${TARGETARCH} /bin/no-response
RUN chmod +x /bin/no-response

####################################################################################################
# no-response
####################################################################################################
FROM scratch AS no-response
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/no-response /bin/no-response
ENTRYPOINT [ "/bin/no-response" ]





