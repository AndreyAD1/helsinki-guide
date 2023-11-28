#! /bin/bash

set -xue

database_url=$1
bot_token=$2
tag=$3
metrics_user=$4
metrics_password=$5
metrics_port=$6

container_name=helsinki-guide
image_name=andreyad/helsinki-guide:$tag

running_container=$(docker ps --all --filter name=$container_name -q)
if [ -n "$running_container" ]
    then
        docker stop $container_name
        docker rm -v $container_name
fi

docker pull $image_name
docker run \
--env DEBUG=1 \
--env DATABASE_URL=$database_url \
--env BOT_TOKEN=$bot_token \
--env METRICS_USER=$metrics_user \
--env METRICS_PASSWORD=$metrics_password \
--env METRICS_PORT=$metrics_port \
--network host \
--log-opt tag=hguide \
--name $container_name \
--detach \
$image_name
