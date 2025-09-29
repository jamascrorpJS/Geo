ab -n 1000 -c 10 localhost:3000/api/v1/location/20

Server Hostname:        localhost
Server Port:            3000

Document Path:          /api/v1/location/20
Document Length:        137 bytes

Concurrency Level:      10
Time taken for tests:   0.157 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      246000 bytes
HTML transferred:       137000 bytes
Requests per second:    6367.36 [#/sec] (mean)
Time per request:       1.571 [ms] (mean)
Time per request:       0.157 [ms] (mean, across all concurrent requests)
Transfer rate:          1529.66 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:     0    1   0.7      1       4
Waiting:        0    1   0.7      1       4
Total:          1    2   0.7      1       4
WARNING: The median and mean for the total time are not within a normal deviation
        These results are probably not that reliable.

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      2
  75%      2
  80%      2
  90%      3
  95%      3
  98%      4
  99%      4

провел тест

