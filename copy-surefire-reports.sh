#!/bin/sh

mkdir $2
find $1 -name "TEST-*.xml" -exec cp {} $2 \;
