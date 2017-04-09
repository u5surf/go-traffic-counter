#!/bin/sh
nc -U $1 < $2 > /dev/null
