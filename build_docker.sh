#!/bin/bash

version=$1

docker_repo="linshenqi/taskmate"

docker build -t ${docker_repo}:${version} .

docker push ${docker_repo}:${version}
