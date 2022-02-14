#!/bin/bash

echo "Temporarily open the port for the Management Console.  This script is"
echo "synchronous.  Hit Ctrl+c to stop port forwarding"
echo
kubectl port-forward --address 0.0.0.0 svc/verticadb-operator-console 5450
