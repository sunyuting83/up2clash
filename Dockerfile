FROM alpine:latest

RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY clash-linux-amd64 /etc/clash/
COPY config.yaml /etc/clash/
COPY ui.zip /etc/clash/

RUN apk add --no-cache \
 ca-certificates  \
 bash  \
 iptables  \
 bash-doc  \
 unzip \
 bash-completion  \
 rm -rf /var/cache/apk/* && \
 chmod a+x /usr/local/bin/entrypoint.sh
RUN unzip /etc/clash/ui.zip -d /etc/clash/
RUN rm -rf /etc/clash/ui.zip

ENTRYPOINT ["entrypoint.sh"]
CMD ["/bin/sh","/etc/up2c/start.sh","start"]