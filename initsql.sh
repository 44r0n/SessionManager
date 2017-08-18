#!/bin/bash
cd "$(dirname "$0")"
if [ -z "$2" ]; then
  DOCKERIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' sessionmanager)
else
  DOCKERIP=$2
fi

echo "Seting up database on $DOCKERIP..."
mysql -uroot -pmypassword -h $DOCKERIP -P 3306 < data/sessionmanager.sql
if ! [ -z "$1" ]; then
  if [ $1 = "-t" ]; then
    echo "Testing database..."
    cd mytap
    mysql -uroot -pmypassword -h $DOCKERIP -P 3306 < mytap.sql
    cd ..
    mysql -uroot -pmypassword -h $DOCKERIP -P 3306 < data/test.sql
  fi
fi
echo "MySQL commands executed at $DOCKERIP"

if [ ! -d "configuration" ]; then
  mkdir configuration
  echo "Folder configuration created"
fi

if [ -e configuration/configuration.json ]; then
    rm configuration/configuration.json
fi

if [ $DOCKERIP = "127.0.0.1" ]; then
  echo '{
    "ip":"'$DOCKERIP'",
    "port":3306,
    "ConnString":"root:mypassword@/sessionmanager"
  }' >> configuration/configuration.json
  echo "configuration.json file created"
  cat configuration/configuration.json
else
  echo '{
    "ip":"'$DOCKERIP'",
    "port":3306,
    "ConnString":"root:mypassword@('$DOCKERIP':3306)/sessionmanager"
  }' >> configuration/configuration.json
  echo "configuration.json file created"
  cat configuration/configuration.json
fi

