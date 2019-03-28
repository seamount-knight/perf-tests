export MASTER_NAME=10.21.64.18
export MasterSSHPass=tawFOFKs3
/bin/bash -c "go run cmd/clusterloader.go --kubeconfig=/root/.kube/config --testconfig=./testing/load/config.yaml --provider=local --masterip=10.21.64.18  --report-dir=./report-2/report-load --logtostderr=true"
