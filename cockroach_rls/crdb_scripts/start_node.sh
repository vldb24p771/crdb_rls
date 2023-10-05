
NETWORK_NAME="rls_network"
CONTAINER_NAME_PREFIX="roach"
MAX_REGIONS=8
NETWORK_BASE="10.62.0."
CONTAINERS_FILE="./containers.yml"
PORT3=`expr $i + 26357`

for ((i=1; i<=$MAX_REGIONS; i++))
do
{
    join=$join$CONTAINER_NAME_PREFIX$i$":"$PORT3$","
}
done
join="${join%?}"

echo $join
for ((i=0; i<$MAX_REGIONS; i++))
do
{
    # each region takes a subnet
    PORT1=`expr $i + 8080`
    PORT2=`expr $i + 26257`
    NETWORK=`expr $i \* 1 + 200`

    IP=$NETWORK_BASE$NETWORK
    NUMBER=`expr $i + 1`

    container_name=$CONTAINER_NAME_PREFIX$NUMBER

    ssh_command=" /root/crdb_origin/cockroach start --store=node$NUMBER --advertise-addr=$container_name:$PORT3 --http-addr=$container_name:$PORT1 --listen-addr=$container_name:$PORT3 --sql-addr=$container_name:$PORT2  --insecure --join=$join --background"
    echo $ssh_command
    echo $IP
    ssh -o StrictHostKeyChecking=no $IP $ssh_command

}&
done
wait
sleep 2