# analyzer-d4-balboa
Ingests Type 8 Passive DNS and writes into a linux socket for balboa https://github.com/DCSO/balboa to consume

# Installation
```bash
go get https://github.com/D4-project/analyzer-d4-balboa
```

# Configuration files
 - balboa_socket: path to the UNIX socket
 - redis: path to the d4 redis server
 - redis_queue: uuid of the analyzer's redis queue
 
 # Use
 ```bash
 $analyzer-d4-balboa -c conf.sample
 ```
 
 # Query Balboa to test
 Once you launched the analyzer, pick one of the domains listed in its output and query Balboa (serving here on http://127.0.0.1:8080):
 ```bash
 #!/bin/bash
curl \
 -X POST \
 -H 'Content-Type: application/json' \
 --data '{"query" : "query{ entries(rrname: \"www.cnn.com\", limit: 1) { rrname rrtype rdata time_first time_last sensor_id count } } "}' http://127.0.0.1:8080/
 ```
