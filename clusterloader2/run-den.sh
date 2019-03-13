export MASTER_NAME=10.21.64.18
export MasterSSHPass=tawFOFKs3
/bin/bash -c "go run cmd/clusterloader.go --kubeconfig=/root/.kube/config --testconfig=./testing/density/config.yaml --provider=kubemark --masterip=10.21.64.18  --report-dir=./report-den --logtostderr=true"
