#!/bin/bash
cd "$(dirname "$0")"
DOCKERIP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' sessionmanager)
echo "Seting up database..."
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

if [ -e configuration/configuration.json ]; then
    rm configuration/configuration.json
fi
echo '{
  "port":3306,
  "ConnString":"root:mypassword@('$DOCKERIP':3306)/sessionmanager"
}' >> configuration/configuration.json
echo "configuration.json file created"
