This is a tutorial to use orbiter as autoscaler with Docker Swarm.

Requirements:

* Docker CLI installed
* Docker Machine
* Virtualbox

1. Create a swarm cluster with 3 nodes. One manager and two workers. It uses
   Docker Machine and Virtualbox as provider.
```
make init
```

2. Point Docker CLI to the right Docker machine. The master is called sw1.

```
eval $(docker-machine env sw1)
```

3. Deploy your stack. The stack contains two services, one called orbiter and
   another called micro. Micro is the web app that we are autoscaling.

```
make deploy
```

4. This command shows the current situation. The first table represent the
   number of services and how many task are running under a service. You should
   see two services. Micro has 3 tasks.

```
make ps
```

5. There is an utility in the makefile to scale up and down micro.

```
make scale-up
make scale-down
```

6. You can check the number of tasks changing every time you scale with the
   command:

```
make ps
```

7. At the end of the test you can clean your environment

```
make destroy
```

That's it. Have a look at the code in `./stack.yml` and `Makefile` to have more
information about how orbiter works.
