export MASTER_NODES=10.0.6.222,10.0.6.51,10.0.6.84
export LOCAL_SSH_KEY=/Users/xuehaishan/Workspace/ssh-key/new
./clusterloader --logtostderr -v 9 --kubeconfig=./config --masterip=10.0.6.222,10.0.6.51,10.0.6.84  --provider=local --report-dir=./report --testconfig=./testing/load/config.yaml
 