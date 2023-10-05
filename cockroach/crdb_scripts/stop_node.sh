for ((i=150; i<250; i=i+10))
do
    echo "Stop the progress in client $i(s)"

    cmd="ps -ef | grep tapir | grep -v grep | awk '{print \$2}' | xargs kill -9"
    echo $cmd
    ssh 10.62.0.$i "ps -ef | grep crdb | grep -v grep | awk '{print \$2}' | xargs kill -9"
    ssh 10.62.0.$i "ps -ef | grep cockroach | grep -v grep | awk '{print \$2}' | xargs kill -9"
done