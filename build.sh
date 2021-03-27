#!/usr/bin/env bash

for SRC in *.go
do
	echo ${SRC}
	BASE=${SRC%.*}
	GOOS=linux GOARCH=amd64 go build -o bin/${BASE}.linux.amd64 ${SRC}
	GOOS=freebsd GOARCH=amd64 go build -o bin/${BASE}.freebsd.amd64 ${SRC}
	GOOS=windows GOARCH=amd64 go build -o bin/${BASE}.exe ${SRC}
done
