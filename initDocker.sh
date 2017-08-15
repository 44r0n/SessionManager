#!/bin/bash
function startdoker {
  echo "Uping docker..."
  docker run --detach --name=sessionmanager --env="MYSQL_ROOT_PASSWORD=mypassword" mysql
  docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' sessionmanager
  
if [ -z "$1" ]; then
  DOCKERIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' sessionmanager)
else
  DOCKERIP=$1
fi

  until nc -z -v -w30 $DOCKERIP 3306
  do
    echo "Waiting for database connection..."
    sleep 5
  done
}

DOKERSTATUS=$(docker inspect -f {{.State.Running}} sessionmanager)

if [ -z $DOKERSTATUS ]; then
  startdoker
else
  if [ $DOKERSTATUS = false ]; then
    docker rm sessionmanager
    startdoker
  fi
fi

if [ -z "$1" ]; then
  DOCKERIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' sessionmanager)
else
  DOCKERIP=$1
fi

echo "Service up and runing at $DOCKERIP"
