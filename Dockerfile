# syntax=docker/dockerfile:1

FROM turbot/steampipe:0.15.0 as builder
ARG GO_DOWNLOAD_LINK=https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
ARG STEAMPIPE_PLUGIN_RSPACE=steampipe-plugin-rspace.plugin

### install Go and build the plugin
USER root
RUN wget -O go.tar.gz $GO_DOWNLOAD_LINK
RUN  tar -C /usr/local -xzf go.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"
RUN go version
WORKDIR /app
COPY go.mod go.sum main.go  ./
ADD rspace ./rspace
RUN go mod download
RUN go build -o /$STEAMPIPE_PLUGIN_RSPACE

FROM turbot/steampipe:0.15.0

ARG STEAMPIPE_PLUGIN_RSPACE_INSTALL=/home/steampipe/.steampipe/plugins/local/rspace
ARG STEAMPIPE_PLUGIN_RSPACE=steampipe-plugin-rspace.plugin
WORKDIR /app
USER steampipe
COPY rspace.spc /home/steampipe/.steampipe/config/
RUN mkdir -p  ${STEAMPIPE_PLUGIN_RSPACE_INSTALL}
USER root
COPY --from=builder --chown=steampipe:0 /$STEAMPIPE_PLUGIN_RSPACE ${STEAMPIPE_PLUGIN_RSPACE_INSTALL}/

## checkout the dashboard
RUN apt-get update && apt-get -y install git vim nano
WORKDIR /git
RUN git clone --depth 1 -bv0.0.1 https://github.com/richarda23/steampipe-mod-rspace.git
RUN chown -R steampipe /git

USER steampipe
WORKDIR /git/steampipe-mod-rspace

CMD ["steampipe", "service", "start", "--foreground", "--dashboard", "--dashboard-listen=network"]





