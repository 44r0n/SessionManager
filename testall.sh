#!/bin/bash
./initDocker.sh $1
./initsql.sh -t $1
./testapp.sh
./initsql.sh "" $1
