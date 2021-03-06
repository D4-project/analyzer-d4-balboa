# Archiving analyzer-d4-balboa
This repository is now archived - to interface D4 with Balboa, the prefered tool is [d4-core generic unix socket exporter](https://github.com/D4-project/d4-core/blob/master/server/analyzer/analyzer-d4-export/d4_export_unix.py).


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
 
# Send PassiveDNS data to d4
```bash
# passivedns -i eth0 -l /dev/stdout | d4-amd64l -c conf.d4server
```

 # Query Balboa GraphQL server
 Once you launched the analyzer, pick one of the domains listed in its output and query Balboa (serving here on http://127.0.0.1:8080):
 ```bash
 #!/bin/bash
curl \
 -X POST \
 -H 'Content-Type: application/json' \
 --data '{"query" : "query{ entries(rrname: \"www.cnn.com\", limit: 1) { rrname rrtype rdata time_first time_last sensor_id count } } "}' http://127.0.0.1:8080/
 ```
