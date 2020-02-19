#!/bin/bash

for SRC in *.go
do
	echo ${SRC}
	GOOS=linux GOARCH=amd64 go build ${SRC}
	GOOS=windows GOARCH=amd64 go build ${SRC}
done
