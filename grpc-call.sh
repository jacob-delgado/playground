#!/bin/bash

grpcurl -plaintext -d '{"name":"test123"}' localhost:8000 playground.v1.PlaygroundService/GetFeature
