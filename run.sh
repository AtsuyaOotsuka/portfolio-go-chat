#!/bin/bash

if ! docker network ls | grep -q "atylab_net"; then
    docker network create \
        --driver bridge \
        --subnet 172.30.0.0/16 \
        --gateway 172.30.0.1 \
        atylab_net
fi

docker compose up --build -d
