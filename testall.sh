#!/bin/bash
./initDocker.sh
./initsql.sh -t
./testapp.sh
./initsql.sh
