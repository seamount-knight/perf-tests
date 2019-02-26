./benchmark --endpoints=http://10.21.128.13:2379 --target-leader --conns=1 --clients=1 put --key-size=8 --sequential-keys --total=10000 --val-size=256

./benchmark --endpoints=http://10.21.128.13:2379 --target-leader --conns=100 --clients=1000 put --key-size=8 --sequential-keys --total=100000 --val-size=256


./benchmark --endpoints=http://10.21.128.6:2379,http://10.21.128.13:2379,http://10.21.128.3:2379 --target-leader --conns=100 --clients=1000  put --key-size=8 --sequential-keys --total=100000 --val-size=256


读取

./benchmark --endpoints=http://10.21.128.6:2379,http://10.21.128.13:2379,http://10.21.128.3:2379 --conns=1 --clients=1 range YOUR_KEY --consistency=l --total=10000


./benchmark --endpoints=http://10.21.128.6:2379,http://10.21.128.13:2379,http://10.21.128.3:2379 --conns=1 --clients=1 range YOUR_KEY --consistency=s --total=10000


./benchmark --endpoints=http://10.21.128.6:2379,http://10.21.128.13:2379,http://10.21.128.3:2379 --conns=100 --clients=1000 range YOUR_KEY --consistency=l --total=100000


./benchmark --endpoints=http://10.21.128.6:2379,http://10.21.128.13:2379,http://10.21.128.3:2379 --conns=100 --clients=1000 range YOUR_KEY --consistency=s --total=100000

     