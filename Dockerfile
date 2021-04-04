FROM alpine:latest

COPY clash /usr/local/bin/
COPY Country.mmdb /clash_config/clash/
COPY entrypoint.sh /usr/local/bin/
COPY config.yaml /clash_config/clash/

RUN apk add --no-cache \
 ca-certificates  \
 bash  \
 iptables  \
 bash-doc  \
 bash-completion  \
 rm -rf /var/cache/apk/* && \
 chmod a+x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
CMD ["/usr/local/bin/clash","-d","/clash_config"]