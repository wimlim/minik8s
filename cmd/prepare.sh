#!/bin/bash

docker network create app-tier --driver bridge

# Start etcd server
docker run -d --name etcd-server --network app-tier --publish 2379:2379 --publish 2380:2380 --env ALLOW_NONE_AUTHENTICATION=yes --env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 bitnami/etcd:latest

# Start RabbitMQ server
docker run --name rabbitmq -d -p 15672:15672 -p 5672:5672 -v /home/rabbitmq/data:/data -e RABBITMQ_DEFAULT_USER=ling -e RABBITMQ_DEFAULT_PASS=123456 rabbitmq:management