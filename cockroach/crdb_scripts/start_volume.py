import subprocess
import yaml


NETWORK_NAME = "rls_network"
CONTAINER_NAME_PREFIX = "roach"
HOSTS_FILE = '/home/ruijie/share/crdb/cockroach/crdb_scripts/hosts.yml'
MAX_REGIONS = 10
NETWORK_BASE = "10.62.0"
IMAGE_NAME = "rjgong/crdb:v2"
CONTAINERS_FILE = "./containers.yml"

hosts = [] 
with open (HOSTS_FILE, 'r') as f:
    data = yaml.load(f)
    hosts = data['hosts']

containers = []

for i in range(MAX_REGIONS):
    # each region takes a subnet
    host = hosts[i % len(hosts)]
    IP = NETWORK_BASE + "." + str(i*10+150)
    containers.append(IP)
    print(host, IP)
    container_name = CONTAINER_NAME_PREFIX + str(i)
    ssh_command = " docker volume create --driver local " + \
                  " --opt type=nfs4 " + \
                  " --opt o=addr=10.22.1.5,rw " + \
                  " --opt device=:/share " + container_name
    print(ssh_command)
    subprocess.call(['ssh', host, ssh_command]) 

with open(CONTAINERS_FILE, 'w+') as f:
    data = {}
    data['containers'] = containers
    yaml.dump(data, f)
