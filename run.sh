#!/usr/bin/env bash

IMAGE=docker.jw4.us/logsrv:latest
NAME=logsrv
PORT=${PORT:-11181}
VERBOSE=${VERBOSE:-}

docker pull ${IMAGE}
docker stop ${NAME}
docker logs ${NAME} &> $(TZ=UTC date +%Y-%m-%d-%H%M-${NAME}.log)
docker rm -v -f ${NAME}

docker run -d \
  --name ${NAME} \
  --restart=always \
  -e VERBOSE="${VERBOSE}" \
  -e PORT="${PORT}" \
  -p ${PORT}:${PORT} \
  ${IMAGE}
