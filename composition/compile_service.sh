#!/usr/bin/env bash

cd ../
./build_service.sh
cd composition
./run_composer.sh up --build