#!/bin/bash

GREEN='\033[0;32m'

printf "${GREEN} Running my analytics service...\n"

cd analytics


cargo build

cargo run
