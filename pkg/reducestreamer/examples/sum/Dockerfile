####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/sum-example /bin/sum-example
RUN chmod +x /bin/sum-example

####################################################################################################
# sum
####################################################################################################
FROM scratch as sum
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/sum-example /bin/sum-example
ENTRYPOINT [ "/bin/sum-example" ]
