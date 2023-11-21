#! /bin/bash

set -xue

database_url=$1
bot_token=$2

container_name=helsinki-guide
image_name=andreyad/helsinki-guide

running_container=$(docker ps --all --filter name=$container_name -q)
if [ -n "$running_container" ]
    then
        docker stop $container_name
        docker rm -v $container_name
fi

docker pull $image_name
docker run \
--env Debug=1 \
--env DatabaseURL=$database_url \
--env BotAPIToken=$bot_token \
--network host \
--log-opt tag=hguide \
--name $container_name \
--detach \
$image_name
