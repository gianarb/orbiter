FROM golang:1.8.3 as builder

WORKDIR /go/src/github.com/gianarb/orbiter
ADD . /go/src/github.com/gianarb/orbiter/
RUN make build

FROM scratch

COPY --from=builder /go/src/github.com/gianarb/orbiter/bin/orbiter /bin/orbiter
# ENTRYPOINT ["orbiter"]
CMD ["/bin/orbiter", "daemon"]
