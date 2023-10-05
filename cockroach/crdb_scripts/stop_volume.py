import yaml
import subprocess

NETWORK_NAME = "rls_network"
CONTAINER_NAME_PREFIX = "crdb_"
HOSTS_FILE = '/home/ruijie/share/test/scripts/hosts.yml'
MAX_REGIONS = 10
NETWORK_BASE = "10.62.0"
IMAGE_NAME = "cockroachdb/cockroach"
CONTAINERS_FILE = "containers.yml"

hosts = [] 
with open (HOSTS_FILE, 'r') as f:
    data = yaml.load(f)
    hosts = data['hosts']

for h in hosts:
    print("stopping volume on " + h)
    ssh_command = "docker volume rm $(docker volume ls -qf dangling=true)"
    print(ssh_command)
    subprocess.call(['ssh', h, ssh_command])

