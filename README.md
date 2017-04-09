# golang-sample-echo-server

## Description 
Byte counter to transffering huge file with unix domain socket.

## Demo
$go build -o server

$./server 

copy unix-domain-socket-file from(e.g. /tmp/hogehoge...)

other terminal

$dd if=/dev/urandom of=sample.dat bs=1G count=1

$./client unix-domain-socket-file sample.dat

## Requirement
go > 1.7

## Install 
$go build -o server server.go


