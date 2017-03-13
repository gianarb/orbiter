FROM debian:sid

ADD ./bin/orbiter /opt/orbiter

CMD ["/opt/orbiter", "daemon", "-config", "/etc/orbiter.yml"]
