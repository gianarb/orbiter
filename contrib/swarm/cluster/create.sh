#!/bin/bash
provider=virtualbox
echo "### Init Servers ###"

docker-machine create -d ${provider} sw1 &
docker-machine create -d ${provider} sw2 &
docker-machine create -d ${provider} sw3 &

wait

echo "### Configurate cluster ###"

MANAGER_IP=$(docker-machine ip sw1)
docker-machine ssh sw1 docker swarm init --advertise-addr ${MANAGER_IP}
TOKEN=$(docker-machine ssh sw1 docker swarm join-token -q worker)
docker-machine ssh sw2 docker swarm join --token ${TOKEN} ${MANAGER_IP}:2377
docker-machine ssh sw3 docker swarm join --token ${TOKEN} ${MANAGER_IP}:2377
docker-machine ssh sw4 docker swarm join --token ${TOKEN} ${MANAGER_IP}:2377

# Information
echo ""
echo "CLUSTER INFORMATION"
echo "discovery token: ${TOKEN}"
echo "Environment variables to connect trough docker cli"
docker-machine env sw1
