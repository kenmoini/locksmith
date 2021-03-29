FROM quay.io/kenmoini/golang-ubi:latest AS build

WORKDIR /opt/app-root/src
COPY . .
RUN make build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest AS bin

COPY --from=build /opt/app-root/src/dist/locksmith /opt/app-root/bin/
COPY container_root/ /

EXPOSE 8080

CMD [ "/opt/app-root/bin/container_start.sh" ]
