
for i in 100 200 300 400 500 600 700 800 900 1000
do
for j in {1..2}
do
echo "kv concurrency is $i for $j"

./cockroach workload run kv --duration=1m 'postgresql://root@roach1:26257?sslmode=disable' --concurrency $i

done
done