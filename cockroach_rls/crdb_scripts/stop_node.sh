for ((i=150; i<250; i=i+10))
do
    echo "Stop the progress in client $i(s)"

    ssh 10.62.0.$i "ps -ef | grep java | grep -v grep | awk '{print \$2}' | xargs kill -9"
done

#ps -ef | grep java | grep -v grep | awk '{print $2}' | sudo xargs kill -9
#ps -ef | grep java | grep -v grep | awk '{print $2}' | sudo xargs kill -9