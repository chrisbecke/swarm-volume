FROM golang:1.19 AS build

WORKDIR /src
COPY . .
WORKDIR /src/glusterfs-swarm-volume
RUN go build

WORKDIR /tini
# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini/tini
RUN chmod +x /tini/tini

FROM gluster/glusterfs-client:latest AS plugin

COPY --from=build /tini/tini /usr/local/bin/tini
COPY --from=build /src/glusterfs-swarm-volume/docker-volume-plugin /usr/local/bin/docker-volume-plugin

ENTRYPOINT ["tini", "--"]
CMD ["docker-volume-plugin"]
