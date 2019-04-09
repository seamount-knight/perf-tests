export MASTER_NAME=10.21.32.13
export MasterSSHPass=tawFOFKs3
/bin/bash -c "go run cmd/clusterloader.go --kubeconfig=/root/.kube/config --testconfig=./testing/load/config.yaml --provider=local --masterip=10.21.32.13  --report-dir=./report-3/report-load --logtostderr=true"
