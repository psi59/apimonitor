FROM golang:alpine
MAINTAINER psi59@lalaworks.com

WORKDIR "/root/.apimonitor"
COPY ./config/server_config.default.yaml /root/.apimonitor/config/server_config.yaml
COPY ./apimonitor.linux /usr/bin/apimonitor
ENV AM_CONFIG_PATH=/root/.apimonitor/config
CMD ["apimonitor"]