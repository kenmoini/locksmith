FROM registry.access.redhat.com/ubi8/go-toolset:latest AS build

WORKDIR /opt/app-root/src
COPY . .
RUN go build

FROM scratch AS bin

COPY --from=build /opt/app-root/src/locksmith /usr/local/bin/
#COPY container_root/ /
RUN mkdir -p /etc/locksmith

EXPOSE 8080

CMD [ "locksmith -config /etc/locksmith/config.yml" ]