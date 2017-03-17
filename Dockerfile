FROM scratch

ADD ./bin/orbiter /bin/orbiter

CMD ["orbiter", "daemon"]
