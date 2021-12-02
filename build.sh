#!/bin/bash

rm -rf out/
mkdir out/

go build -ldflags "-X main.environment=cli" -o out/papyrus
go build -ldflags "-X main.environment=web" -o out/papyrus-web
