# prometheus

## Introduction: "Hidden Traps" in Monitoring Systems

In morder distributed system, as a leading monitoring and alert tool, prometheus get femous since its data scraptionmg and querying ability.
However, as system complexity get more, "" issue become clear, become the " killer" by impacting performance and expanding. the blog will deeply touch the original, influence and tickle tactics.


the nature of high cardinality label problem
what is high cardinality label problem?

The high cardinality label problemis that in prometheus, some label's value is too heavy. For example, labels like `user_id`, `transactionid`, if use unsuitable, will conduct the accelerated increasing of cardinality label, and increase the performance of calculating when querying.

impacting from the issue?
1. storage pressure increasing: each signle label mix will generate a timeseries, the high cardinality label will conduct the number of timeseries increasing, contine increase storage requirement, and may even overhit the storage ability of prometheus.

2. query performance descreasing: THE ISSUE WIKLl made the querying time nned to query more data when fit labels, which will conduct the querying react time longer, and impact the performance of realtime monitoring and alerting.

3. memory consumption acceleate: when Prometheus was handing high cardinality labels, it need more memory to manage and index these labels, which might conduct the jinzhangxxx of system resouce, and continue to impact the stablebility in zhengtixxx

So, in low level (storage design level), why is the problem so scary?

let's deep into the structure.



what is high cardinality label?
when you record http request status, something like `http_requests_total{method="GET", handler="/api", status="200"}` , if you want to record `user_id`, if there are tons of tons user in your system, the label will bloom.


what is timeseries data?

TODO chart two-dimention, how to draw it?












## 1. what is timeseries?

- data schema
  - identifier -> (t0, v0), (t1, v1), (t2, v2), ...
- prometheus data model
  - <metrics name>{<label name>=<label name>, ...}
- typical set of series identifiers (data model)
  - {    "__name__": "http_requests_total",    "pod": "example-pod",    "job": "example-job",    "path": "/api/v1/resource",    "status": "200",     "method": "GET"} @143xxx 94355
- query
    - __name__="request_total" - 
    - method="PUT|POST" - 
1. ideal R/W


TODO 二维象限 (prometheus 1.x 自己的存储设计是二维象限那样的)，展示，如果label 爆炸 series 会爆炸 -- 毫无疑问

每个series 是一行，这样就太多了--- 多少行，就有多少文件（并行写）
---> 批量处理，因为在一个文件里

那么，问题是什么？ Fundamental Problem

TODO 二维象限
1. storage problem
   1. IDE - spinning phscially
   2. SSD - write amplification
      1. 4KB MODIFY
      2. 256 DELETE
   3. Query is much more complicated than write 
   4. Time series query could cause the random read
2. ideal read - batch
3. ideal write - 顺序 - 联想到磁盘


### Prometheus Solution (v1.x "V2")

每个series 都是一个文件 -> inode

Series Churn

### TSDB

我希望是顺序写，批量读

操作的文件不一样了，按 block 来

#### when you run `tree ./data`:
```bash
data
├── blocks
│   ├── 01FG7C0Q0G3QZ8N1KZ4P1MXY7V
│   │   ├── chunks
│   │   │   ├── 000001
│   │   │   ├── 000002
│   │   │   └── ...
│   │   ├── index
│   │   │   ├── index.db
│   │   │   └── ...
│   │   ├── meta.json
│   │   └── tombstones
│   │       └── tombstones.json
│   └── 01FG7C0Q0G3QZ8N1KZ4P1MXY7W
│       ├── chunks
│       ├── index
│       ├── meta.json
│       └── tombstones
├── chunks_head
│   ├── 000001
│   ├── 000002
│   └── ...
├── index
│   ├── index.db
│   └── ...
├── meta.json
├── lock
└── wal
    ├── 00000001
    ├── 00000002
    └── ...
```

### Fundamental Design - v3

just simplely explain the conceptions:
- blocks/
  - chunks/: store seperated blocks of timeseries data 
  - index/: store index file, like `index.db`, to accelerate query speed
  - meta.json/: meta data files to describe the data block, including time range and other config info.
  - tombstones/: store deleted info about time series, files like `` was record the deleted data
- chunks_head: store unpacked time series data in memory, data in this part has not been persistent into `blocks/`
- index: gloval index 目录xxx，including index files for querying quickly, like `index.db`
- meta.json: global meta data, record status and configs info of whole TSDB.
- lock: file lock, to avoid multiple prometheus instance acces same data folder 同时地, keep the data persistently
- wal(write-ahead log): record latest writing operations, to keep the data will not loss when electric interrupted or crashed. name of log files walways using numbers, such as `00000001`, `00000002`, etc.


Notes:
- the data will be persisted into disk every 2 hours.
- WAL is used for data recovery
- 2 hour s block could make thew range data query efficiently


Let's delve into the details (analysis one by one):
### block
1. a block will store a time block, maybe 2 hours
   ```
    数据块 1：从 2024-01-01 00:00 到 2024-01-01 02:00
    数据块 2：从 2024-01-01 02:00 到 2024-01-01 04:00
   ```
2. a timeseries data has a number id, 1/2/3/4/5

chunk-head

TODO chart

### disk
3. 列式存储和一般的行存储的对比 - https://chatgpt.com/c/672ae026-16f0-8011-adba-f704565085d7

压缩和优化：Prometheus 使用高效的压缩算法（如 Gorilla 压缩）来减少时间序列数据的存储空间。

### WAL

TODO

### index
4. 倒排索引：
   ```
    id=1: [TS1, TS2, TS3]
    id=2: [TS4, TS5]
    status=active: [TS1, TS5, TS6]
   ```

TODO chart


TODO index
how does the index work in Prometheus?
```json
{
    "__name__": "http_requests_total",
    "pod": "example-pod",
    "job": "example-job",
    "path": "/api/v1/resource",
    "status": "200",
    "method": "GET"
}
```
label to filtering
time range to filtering
like k-v, but multiple key pairs.
- inverted index
- unique ID for every series
- conduct the label's index
- in short
  - nuber if labels is significantly lesse then the bnumber of series
  - walking through all of the labels is not problwem


when I filter "status=404", the index will told us: 
`1 2 4 5...`
when we fileter "method=GET", the index will told us:
`2 4 6 7`
and now, the time series will return the merged result: ID=2/4
and follwo the ID, the prometheus will return the real data: 
```json
{
    "__name__": "http_requests_total",
    "pod": "example-pod",
    "job": "example-job",
    "path": "/api/v1/resource",
    "status": "200",
    "method": "GET"
}, {
    "__name__": "http_requests_total",
    "pod": "example-pod-2",
    "job": "example-job-2",
    "path": "/api/v2/resource",
    "status": "404",
    "method": "POST"
}
```


large file with "mmap"
- mmap stands for memory-mapped files. it is a way to read and write files without invoking **system calls**.
- it is great if multiple processes accessing data in a read only fashion from the same file.
- it allows all those processes to share the same physical memory pages, saivntg a lot of memory/
- it also allows the operating system to optimize paging operations

TODO ![](./mmap_prom.png)

how to handle it?
1. first, you can use the xxx as the label
   1. like use index to index the 
      1. `mysql` 中，what items should not as the index?
2. stardalization and leveling:
   1. you should make you label stardalizatioon, managing in diffenrent level to avoid differnent coard mixing in use, keep clean and consistent in label system
3. data pre-processing and merging: before the data get into prometheus, do something must pre-producing and merging, decreasing the explosure of high cardinality label directly. For example, using log tools to merge data before you merge resume to prometheus

SUMMARY
some performance increaement techinques in prometheus you should know:
1. block storage
2. y压缩
   1. timestamp and values
3. index
4. mmap
5. compaction