./cockroach workload init kv --zipfian \
'postgresql://root@localhost:26257?sslmode=disable'

./cockroach workload run kv \
--duration=1m \
'postgresql://root@localhost:26257?sslmode=disable'

