FROM registry.access.redhat.com/ubi8/go-toolset:1.14.12 AS build

WORKDIR /opt/app-root/src
COPY . .
RUN go build

FROM scratch AS bin

COPY --from=build /opt/app-root/src/locksmith /usr/local/bin/

EXPOSE 8080

CMD [ "locksmith -config /etc/locksmith/config.yml" ]