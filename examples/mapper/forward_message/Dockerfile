####################################################################################################
# base
####################################################################################################
FROM alpine:3.20 AS base
ARG TARGETARCH
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/forward-message-example-${TARGETARCH} /bin/forward-message-example
RUN chmod +x /bin/forward-message-example

####################################################################################################
# forward-message
####################################################################################################
FROM scratch AS forward-message
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/forward-message-example /bin/forward-message-example
ENTRYPOINT [ "/bin/forward-message-example" ]
