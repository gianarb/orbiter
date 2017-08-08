IP = $(shell docker-machine ip sw1)
init:
	./cluster/create.sh
destroy:
	docker-machine rm -f sw1 sw2 sw3
deploy:
	docker stack deploy --compose-file stack.yml stack
ps:
	docker service ls
	docker stack ps stack
scale-up:
	curl -X POST -v http://$(IP):8081/v1/orbiter/handle/autoswarm/stack_micro/up
scale-down:
	curl -X POST -v http://$(IP):8081/v1/orbiter/handle/autoswarm/stack_micro/down
