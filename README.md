[![Build
Status](https://travis-ci.org/gianarb/orbiter.svg?branch=master)](https://travis-ci.org/gianarb/orbiter)

Public and private cloud or different technologies like virtual machine or
containers. Our applications and our environments require to be resilient
doesn't matter where they are or which services are you using.

This project is a work in progress cross platform open source autoscaler.

We designed in collaboration with InfluxData to show how metrics can be used.

It is based on plugins called `provider`. At the moment we implemented:

* Docker Swarm mode
* DigitalOcean

```sh
orbiter daemon -c config.yml
```
Orbiter is a daemon that use a YAML configuration file to starts one or more
autoscaler and it exposes an HTTP API to trigger scaling up or down.

```yaml
autoscalers:
  events:
    provider: swarm
    parameters:
    policies:
      php:
        up: 4
        down: 3
  infra_scale:
    provider: digitalocean
    parameters:
      token: zerbzrbzrtnxrtnx
      region: nyc3
      size: 512mb
      image: ubuntu-14-04-x64
      key_id: 163422
      # https://www.digitalocean.com/community/tutorials/an-introduction-to-cloud-config-scripting
      userdata: |
        #cloud-config

        runcmd:
          - sudo apt-get update
          - wget -qO- https://get.docker.com/ | sh
    policies:
      frontend:
        up: 2
        down: 3
```
This is an example of configuration file. Right now we are creating two
autoscalers one to deal with swarm called `events/php` and the second one with
DigitalOcean called `infra_scale`.

## Vocabulary

* `provider` contains integration with the platform to scale. It can be swarm,
  DigitalOcean, OpenStack, AWS.
* Every autoscaler supports a k\v storage of `parameters` that you can use to
  configure your provider.
* autoscaler` is composed by provider, parameters and policies. You can have
  one or more.
* autoscaler has only or more policies they contains information about a
  specific application.

You can have only one autoscaler or more with th same provider. Same for
policies, only one or more. Doesn't matter.

## Http API
Orbiter exposes an HTTP JSON api that you can use to trigger scaling UP (true)
or DOWN (false).

The concept is very simple, when your monitoring system knows that it's time to
scale it can call the outscaler to persist the right action

```sh
curl -v -d '{"direction": true}' \
    http://localhost:8000/handle/infra_scale/docker
```

## Embeddable
This project is trying to provide also an easy API to maintain a lot complexy
and clean code base in order to allow you to use `orbiter` as project for your
applications.
OpenStack, Kubernets all of them have a sort of autoscaling feature that you can
use. The idea is to keep this complexy out of your deployment tools. You can
just implement `orbiter`.
Another use case is a self-deployed application. [Kelsey
Hightower](https://www.youtube.com/watch?v=nhmAyZNlECw) had a talk about this
idea. I am still not sure that can be real in those terms but we are already
moving something in our applications that before was in external system as
monitoring, healthcheck why not the deployment part?

```go
package scalingallthethings

import (
	"github.com/gianarb/orbiter/autoscaler"
	"github.com/gianarb/orbiter/provider"
)

func CreateAutoScaler() *autoscaler.Autoscaler{
    p, _ := provider.NewProvider("swarm", map[string]string{})
    a, _ := autoscaler.NewAutoscaler(p, "name-service", 4, 3)
    return a
}
```
