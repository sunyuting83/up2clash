FROM alpine:latest

RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY clash-linux-amd64 /etc/clash/
COPY Country.mmdb /etc/clash/
COPY config.yaml /etc/clash/
COPY ui.zip /etc/clash/
COPY entrypoint.sh /usr/bin/
COPY clash.sh /etc/clash/
COPY up2c /etc/up2c/
COPY start.sh /etc/up2c/

RUN apk add --no-cache \
 ca-certificates  \
 bash  \
 iptables  \
 bash-doc  \
 curl \
 bash-completion
RUN rm -rf /var/cache/apk/*
RUN chmod a+x /usr/bin/entrypoint.sh
RUN chmod a+x /etc/clash/clash.sh
RUN chmod a+x /etc/up2c/start.sh
RUN chmod a+x /etc/clash/clash-linux-amd64
RUN chmod a+x /etc/up2c/up2c
RUN unzip /etc/clash/ui.zip -d /etc/clash/
RUN rm -rf /etc/clash/ui.zip

ENTRYPOINT ["/usr/bin/entrypoint.sh"]