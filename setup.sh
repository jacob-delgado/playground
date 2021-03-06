#!/bin/bash

minikube start --profile custom
skaffold config set --global local-cluster true
eval "$(minikube -p custom docker-env)"
