#!/bin/sh
for i in `seq 1 10` ; do
  nc -U $1 < $2 > /dev/null &
	echo $i
done
