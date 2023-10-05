import yaml
import subprocess

HOSTS_FILE = '/home/ruijie/share/crdb/cockroach/crdb_scripts/hosts.yml'

hosts = [] 
with open (HOSTS_FILE, 'r') as f:
    data = yaml.load(f)
    hosts = data['hosts']

for h in hosts:
    print("stopping container on " + h)
    ssh_command = "docker stop $(docker ps -q)"
    print(ssh_command)
    subprocess.call(['ssh', h, ssh_command])

for h in hosts:
    print("stopping container on " + h)
    ssh_command = "docker rm $(docker ps -aq)"
    print(ssh_command)
    subprocess.call(['ssh', h, ssh_command])
