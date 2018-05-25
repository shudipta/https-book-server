#!/usr/bin/env bash
set -xe

kubectl delete secret -l app=https
kubectl delete service -l app=https
kubectl delete deployment -l app=https
