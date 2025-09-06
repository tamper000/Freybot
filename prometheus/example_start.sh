#!/bin/bash

sudo docker run \
  --name prometheus \
  --network="host" \
  -v $(pwd)/prometheus.yaml:/etc/prometheus/prometheus.yml \
  prom/prometheus
