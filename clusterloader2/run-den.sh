export MASTER_NAME=10.21.128.6
export MasterSSHPass=oue1W2Ks3
/bin/bash -c "go run cmd/clusterloader.go --kubeconfig=/root/.kube/config --testconfig=./testing/density/config.yaml --provider=kubemark --masterip=10.21.128.6  --report-dir=./report --logtostderr=true"
