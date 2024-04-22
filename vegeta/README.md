### Vegeta attack before tuning

```shell
vegeta attack -targets=targets.txt -duration=10s -rate=1000/1s | vegeta report -type=text -output=without_caching.txt
```

Results:

```
BenchmarkHandler_sourceCampaigns-8   	     860	   1214775 ns/op
```

```
Requests      [total, rate, throughput]         10000, 1000.14, 243.86
Duration      [total, attack, wait]             38.248s, 9.999s, 28.249s
Latencies     [min, mean, 50, 90, 95, 99, max]  777.542µs, 1.909s, 1.364ms, 2.627ms, 28.326s, 28.393s, 28.412s
Bytes In      [total, mean]                     3357364, 335.74
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           93.27%
Status Codes  [code:count]                      200:9327  500:673  
Error Set:
500 Internal Server Error
```

### Vegeta attack after tuning
```shell
vegeta attack -targets=targets.txt -duration=10s -rate=1000/1s | vegeta report -type=text -output=without_caching.txt
```
Results

```
BenchmarkHandler_sourceCampaigns-8   	 7146126	       163.5 ns/op
```
```
Requests      [total, rate, throughput]         200000, 20000.02, 19999.93
Duration      [total, attack, wait]             10s, 10s, 45.666µs
Latencies     [min, mean, 50, 90, 95, 99, max]  28.292µs, 4.012ms, 43.439µs, 1.036ms, 8.251ms, 167.958ms, 355.446ms
Bytes In      [total, mean]                     71390000, 356.95
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:200000  
Error Set:
```

