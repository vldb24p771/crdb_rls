import subprocess
import yaml


NETWORK_NAME = "rls_network"
CONTAINER_NAME_PREFIX = "roach"
HOSTS_FILE = '/home/ruijie/share/cockroach/crdb_scripts/hosts.yml'
MAX_REGIONS = 8
NETWORK_BASE = "10.62.0"
IMAGE_NAME = "rjgong/crdb:v2"
CONTAINERS_FILE = "./containers.yml"

hosts = [] 
with open (HOSTS_FILE, 'r') as f:
    data = yaml.load(f,Loader=yaml.FullLoader)
    hosts = data['hosts']

containers = []

for i in range(MAX_REGIONS):
    # each region takes a subnet
    host = hosts[i % len(hosts)]
    PORT1 = str(i+8080)
    PORT2 = str(i+26357)
    IP = NETWORK_BASE + "." + str(i*1+200)
    containers.append(IP)
    print(host, IP)
    print(PORT1, PORT2)
    container_name = CONTAINER_NAME_PREFIX + str(i+1)
    ssh_command = "docker run -t -d -v /home/ruijie/share/cockroach:/root/cockroach -v /home/ruijie/share/crdb_origin:/root/crdb_origin --cap-add NET_ADMIN --network " + NETWORK_NAME + " -p " + PORT1 + ":" + PORT1 + " -p " + PORT2 + ":" + PORT2 + " --name " + container_name + " --ip " + IP + " " + IMAGE_NAME + " /root/cockroach/crdb_scripts/tc_ssh.sh " + IP
    print(ssh_command)
    subprocess.call(['ssh', host, ssh_command])

with open(CONTAINERS_FILE, 'w+') as f:
    data = {}
    data['containers'] = containers
    yaml.dump(data, f)
'''
 + " && " + " /root/crdb/cockroach/cockroach start " +  " --advertise-addr=" + container_name + ":26358" + " --http-addr=" + container_name + ":" + PORT1 + " --listen-addr=" + container_name + ":26357" + " --sql-addr=" + container_name + ":" + PORT2 + " --insecure" + " --join=roach1:26357,roach2:26357,roach3:26357 "
    ssh_command = " docker run -d" + " --name=" + container_name + \
                  " --hostname=" + container_name + \
                  " --net=" + NETWORK_NAME + \
                  " -p " + PORT2 + ":" + PORT2 + \ 
                  " -p " + PORT1 + ":" + PORT1 + \
                  " -v " + container_name + ":/test" + \
                  " cockroachdb/cockroach start " + \
                  " --advertise-addr=" + container_name + ":26357" + \
                  " --http-addr=" + container_name + ":" + PORT1 + \
                  " --listen-addr=" + container_name + ":26357" + \
                  " --sql-addr=" + container_name + ":" + PORT2 + \
                  " --insecure" + \
                  " --join=roach1:26357,roach2:26357,roach3:26357 "
" docker run -d \ " + 
" --name " + container_name + " \ " +
" --network " + NETWORK_NAME + " \ " +
" -p " + "26257:26257" + " \ " +
" -p " + "8080:8080" + " \ " +
" -v " + "/home/ruijie/share/test:/tapir-master" + " \ " +
" rjgong/crdb cd crdb && chmod +x bazelisk-linux-amd64 && cp bazelisk-linux-amd64 /usr/local/bin/bazel && ln -s /usr/local/bin/bazel/bazelisk-linux-amd64 /usr/bin/bazel && ./crdb/cockroach/cockroach start \
  --advertise-addr=roach1:26357 \
  --http-addr=roach1:8080 \
  --listen-addr=roach1:26357 \
  --sql-addr=roach1:26257 \
  --insecure \
  --join=roach1:26357,roach2:26357,roach3:26357
"

" docker volume create --driver local \
  --opt type=nfs4 \
  --opt o=addr=10.22.1.5,rw \
  --opt device=:/share \
  roach " + str(i)

docker run -t -d -v /home/ruijie/share/crdb/cockroach:/root/cockroach -v /home/ruijie/share/cockroach:/root/crdb_origin --cap-add NET_ADMIN --network rls_network --name roach10 --ip 10.62.0.190 rjgong/crdb:v2 /root/cockroach/crdb_scripts/tc_ssh.sh 10.62.0.190


docker run -d \
--name=roach1 \
--hostname=roach1 \
-p 26257:26257 \
-p 8080:8080 \
-v /home/ruijie/share/crdb:/root/crdb -v /home/ruijie/share/cockroach:/root/crdb_origin \
/root/cockroach/cockroach start \
  --advertise-addr=roach1:26357 \
  --http-addr=roach1:8080 \
  --listen-addr=roach1:26357 \
  --insecure \
  --join=roach1:26357,roach2:26357,roach3:26357 \
  --background


  /root/cockroach/cockroach start \
    --advertise-addr=roach2:26357 \
    --http-addr=roach2:8081 \
    --listen-addr=roach2:26357 \
    --insecure \
    --join=roach1:26357,roach2:26357,roach3:26357


  /root/cockroach/cockroach start \
    --store=node3 \
    --advertise-addr=roach3:26357 \
    --http-addr=roach3:8082 \
    --listen-addr=roach3:26357 \
    --insecure \
    --join=roach1:26357,roach2:26357,roach3:26357


./cockroach start \
  --advertise-addr=roach1:26357 \
  --http-addr=roach1:8080 \
  --listen-addr=roach1:26357 \
  --insecure \
  --sql-addr=roach1:26257 \
  --join=roach1:26357,roach2:26357,roach3:26357 \
  --background

./crdb_scripts/start_node.sh

./cockroach --host=roach1:26357 init --insecure

./cockroach workload init kv \
'postgresql://root@roach1:26257?sslmode=disable' \
--insert-count 3000000

cp /root/.cache/bazel/_bazel_root/7dbb1e4f7d0bee92cebfd0a87469e4d9/execroot/com_github_cockroachdb_cockroach/bazel-out/k8-fastbuild/testlogs/pkg/kv/kvserver/tscache/tscache_test/test.log . && chmod 777 test.log

bazel test pkg/kv/kvserver/tscache:all

cp ~/node1/logs/cockroach.log . && chmod 777 cockroach.log

ps -ef | grep bazel | grep -v grep | awk '{print $2}' | xargs kill -9

./cockroach workload run kv \
--duration=1m \
'postgresql://root@roach1:26257?sslmode=disable' \
--concurrency 1000

chmod +x crdb/bazelisk-linux-amd64 && cp crdb/bazelisk-linux-amd64 /usr/local/bin/bazel && ln -s /usr/local/bin/bazel/bazelisk-linux-amd64 /usr/bin/bazel && cd cockroach && bazel build pkg/cmd/cockroach && rm cockroach && cp _bazel/bin/pkg/cmd/cockroach/cockroach_/cockroach .

bazel build pkg/cmd/cockroach && rm cockroach && cp _bazel/bin/pkg/cmd/cockroach/cockroach_/cockroach .


ps -ef | grep java | grep -v grep | awk '{print $2}' | sudo xargs kill -9

rm cockroach && cp _bazel/bin/pkg/cmd/cockroach/cockroach_/cockroach .


./cockroach workload init tpcc \
'postgresql://root@roach1:26257?sslmode=disable'

./cockroach workload run tpcc \
--warehouses=2500 \
--duration=1m \
'postgresql://root@roach1:26257?sslmode=disable' \
--concurrency 100


./cockroach workload init ycsb \
'postgresql://root@roach1:26257?sslmode=disable' \
--insert-count 3000000

./cockroach workload run ycsb \
--duration=1m \
'postgresql://root@roach1:26257?sslmode=disable' \
--concurrency 100


./cockroach workload fixtures import tpcc \
--warehouses=250 \
'postgres://root@roach1:26257?sslmode=disable'

postgres://root@roach1:26257?sslmode=disable postgres://root@roach2:26257?sslmode=disable postgres://root@roach3:26257?sslmode=disable
postgres://root@roach4:26257?sslmode=disable postgres://root@roach5:26257?sslmode=disable postgres://root@roach6:26257?sslmode=disable
postgres://root@roach7:26257?sslmode=disable postgres://root@roach8:26257?sslmode=disable postgres://root@roach9:26257?sslmode=disable
postgres://root@roach10:26257?sslmode=disable

./cockroach workload run tpcc \
--warehouses=2500 \
--ramp=1m \
--duration=5m \
$(cat addrs)


/root/crdb_origin/cockroach start --store=node1 --advertise-addr=roach1:26357 --http-addr=roach1:8080 --listen-addr=roach1:26357 --sql-addr=roach1:26257 --insecure 
--join=roach1:26357,roach2:26357,roach3:26357,roach4:26357,roach5:26357,roach6:26357,roach7:26357,roach8:26357,roach9:26357,roach10:26357 --background
10.62.0.100
/root/crdb_origin/cockroach start --store=node2 --advertise-addr=roach2:26357 --http-addr=roach2:8081 --listen-addr=roach2:26357 --sql-addr=roach2:26258 --insecure --join=roach1:26357,roach2:26357,roach3:26357,roach4:26357,roach5:26357,roach6:26357,roach7:26357,roach8:26357,roach9:26357,roach10:26357 --background

'''

