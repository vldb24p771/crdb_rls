
for i in 1000 2000 3000
do
for j in {1..2}
do
echo "kv concurrency is $i for $j"

#/cockroach workload run ycsb --duration=1m 'postgresql://root@roach1:26257?sslmode=disable' --concurrency $i

./cockroach workload run kv --duration=1m 'postgresql://root@roach1:26257?sslmode=disable' --concurrency $i
done
done