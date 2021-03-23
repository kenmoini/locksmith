FROM quay.io/kenmoini/golang-ubi:latest AS build

WORKDIR /opt/app-root/src
COPY . .
RUN make build

FROM scratch AS bin

COPY --from=build /opt/app-root/src/dist/locksmith /usr/local/bin/
#COPY container_root/ /
#RUN mkdir -p /etc/locksmith

EXPOSE 8080

CMD [ "locksmith -config /etc/locksmith/config.yml" ]