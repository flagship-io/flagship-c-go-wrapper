#!/bin/bash

docker build -t builder .
docker run --rm -v $(pwd):/root builder