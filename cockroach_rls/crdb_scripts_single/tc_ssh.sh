#!/bin/bash


/root/cockroach/crdb_scripts_single/setup-tc.sh $1

echo "start ssh server"
echo "MaxSessions 100" >> /etc/ssh/sshd_config
echo "MaxStartups 100" >> /etc/ssh/sshd_config

/usr/sbin/sshd -D
