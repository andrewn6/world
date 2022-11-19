#!/bin/bash

GREEN='\033[0;32m'

printf "${GREEN} Building my CDN..\n"

cd cdn

go run main.go
