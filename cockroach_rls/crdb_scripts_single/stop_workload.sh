for ((i=100; i<160; i=i+10))
do
    echo "Stop the workload in client $i(s)"

    ssh 10.62.0.$i "ps -ef | grep bazel | grep -v grep | awk '{print \$2}' | xargs kill -9"
    ssh 10.62.0.$i "ps -ef | grep workload | grep -v grep | awk '{print \$2}' | xargs kill -9"
done