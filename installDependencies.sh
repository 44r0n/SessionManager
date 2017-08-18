#!/bin/bash
echo "Installing dependencies..."
go get -v github.com/julienschmidt/httprouter
go get -v github.com/go-sql-driver/mysql
go get -v github.com/tkanos/gonfig
go get -v github.com/dgrijalva/jwt-go
go get -v golang.org/x/crypto/bcrypt
go get -v github.com/smartystreets/goconvey/convey