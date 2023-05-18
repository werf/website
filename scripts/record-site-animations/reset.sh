#!/bin/bash

k3d registry delete --all || true
k3d cluster delete || true

k3d registry create registry.sample-app.test --port 5000 --image ghcr.io/werf/test-registry:latest
k3d cluster create --api-port 6550 -p "80:80@loadbalancer" --registry-use k3d-registry.sample-app.test:5000
