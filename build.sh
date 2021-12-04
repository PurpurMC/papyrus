#!/bin/bash

rm -rf out/
mkdir out/

go build -o out/papyrus
cp web/docs out/ -r
