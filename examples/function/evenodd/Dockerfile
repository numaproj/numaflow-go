####################################################################################################
# base
####################################################################################################
FROM alpine:3.12.3 as base
RUN apk update && apk upgrade && \
    apk add ca-certificates && \
    apk --no-cache add tzdata

COPY dist/go-even-odd-example /bin/go-even-odd-example
RUN chmod +x /bin/go-even-odd-example

####################################################################################################
# evenodd
####################################################################################################
FROM scratch as evenodd
ARG ARCH
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/go-even-odd-example /bin/go-even-odd-example
ENTRYPOINT [ "/bin/go-even-odd-example" ]
