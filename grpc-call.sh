#!/bin/bash

# grpc health checking
grpcurl -d '{service:inventory}' -plaintext localhost:8000 grpc.health.v1.Health/Watch
grpcurl -d '{service:inventory}' -plaintext localhost:8000 grpc.health.v1.Health/Check 

# grpc rpc
grpcurl -d '{name:fishing}' -plaintext localhost:8000 inventory.v1.InventoryService/GetInventory
