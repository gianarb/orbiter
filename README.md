![](https://raw.githubusercontent.com/gianarb/orbiter/master/design/logo.png "Orbiter")


[![Build
Status](https://travis-ci.org/gianarb/orbiter.svg?branch=master)](https://travis-ci.org/gianarb/orbiter)

Orbiter is an easy to run autoscaler for Docker Swarm. It is designed to be
automatic and easy.

We designed in collaboration with InfluxData to show how metrics can be used to
create automation around Docker tasks.


```sh
orbiter daemon
```
Orbiter is a daemon that exposes an HTTP API to trigger scaling up or down.

## Http API
Orbiter exposes an HTTP JSON api that you can use to trigger scaling UP (true)
or DOWN (false).

The concept is very simple, when your monitoring system knows that it's time to
scale it can call the outscaler to persist the right action

```sh
curl -v -d '{"direction": true}' \
    http://localhost:8000/v1/orbiter/handle/infra_scale/docker
```
Or if you prefer

```sh
curl -v -X POST http://localhost:8000/v1/orbiter/handle/infra_scale/docker/up
```

You can look at the list of services managed by orbiter:

```sh
curl -v -X GET http://localhost:8000/v1/orbiter/autoscaler
```

Look at the health to know if everything is working:

```sh
curl -v -X GET http://localhost:8000/v1/orbiter/health
```

## Autodetect

```
orbiter daemon
```

It's going to start in autodetect mode. This modality at the moment only fetches
for Docker SwarmMode. It uses the environment variables DOCKER_HOST (and others) to
locate a Docker daemon. If it's in SwarmMode, orbiter is going to look at all the
services currently running.

If a service is labeled with `orbiter=true` it's going to auto-register the
service and it's going to enable autoscaling.

If a service is labeled with `orbiter=true` orbiter will auto-register the
service providing autoscaling capabilities.

Let's say that you started a service:

```
docker service create --label orbiter=true --name web -p 80:80 nginx
```

When you start orbiter, it's going to auto-register an autoscaler called
`autoswarm/web`. By default up and down are set to 1 but you can override
them with the label `orbiter.up=3` and `orbiter.down=2`.

This scalability allows you to instantiate orbiter in an extremely easy way in
Docker Swarm.

A background job reads the Docker Swarm Event api to keep the services
registered in sync.

## With docker

```sh
docker run -it -v ${PWD}/your.yml:/etc/orbiter.yml -p 8000:8000 gianarb/orbiter daemon
```
We are supporting an image `gianarb/orbiter` in hub.docker.com. You can run it
with your configuration.

In this example I am using volumes but if you have a Docker Swarm 1.13 up and
running you can use secrets.

