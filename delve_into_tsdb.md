

## 1. introduction to Time Series Data

TD;LD
Time series data is a sequence of data points collected or recorded at specific time intervals, typically used to track changes or trends over time.


data schema
  - identifier -> (t0, v0), (t1, v1), (t2, v2), ...

data in prometheus
```
  - <metrics name>{<label name>=<label name>, ...}
```

format it， - typical set of series identifiers (data model)

```json
{ "__name__": "http_requests_total", "pod": "example-pod", "job": "example-job", "path": "/api/v1/resource", "status": "200", "method": "GET"} @1430000000 94355
{ "__name__": "http_requests_total", "pod": "example-pod", "job": "example-job", "path": "/api/v1/resource", "status": "200", "method": "PUT"} @1435000000 94355
{ "__name__": "http_requests_total", "pod": "example-pod", "job": "example-job", "path": "/api/v1/resource", "status": "200", "method": "POST"} @1439999999 94355

```

what parts in it:
- key: series
  - metrics name: __name__
  - labels: 
    ```
    {"pod": "example-pod", "job": "example-job", "path": "/api/v1/resource", "status": "200", "method": "GET"}
    ```
  - timestamp: Timestamp
- value: sample

TODO a chart

How to query:
- query
    - __name__="http_requests_total" - selects all series belonging to the `http_requests_total` metrics
    - method="PUT|POST" - select all series metheod is PUT or POST 


## 2. Challenges in Time Series Storage


write: 


read：

TODO 二维模型图

### The Fundamental Problem
1. storage problem
   1. IDE - spinning phscially
   2. SSD - write amplification
      1. 4KB MODIFY
      2. 256 DELETE
   3. Query is much more complicated than write 
   4. Time series query could cause the random read
2. ideal read - batch
3. ideal write - 顺序 - 联想到磁盘

TODO what is 写放大

TODO how to write/read

## 3. Prometheus Tmwe series Storage solutions


one file per time series:
batch up 1KiB chunks in memory

TODO pic



dark side?
-
-
-
-


issue case in production environment: series churn




## 4. Detailed Componerntes of TSDB im Prometheus

from highest level, just see the file tree:
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


Let's delve into the details
#### Block Structure - little database

tree concept

New Design's Benefits
TODO before copy, think first

从本质上，无非就是读写以及存储压力、performace 的问题：
1. 之前是每个time series 一个file
2. 之前read
   1. 2. 之前如何需要组合怎么办？比如or？（处理两个文件）
3. 之前write
   1. 一个时间点比如100 series 一起，操作100个文件
4. hot spot data 
   1. in memory - recently used
   2. 之前呢？
5. chunk size 



> 思考 --- 计算机世界，直接存raw data，然后每次处理一大波的做法，肯定是没有index + raw data 来的好
好好的设计data model，通过一些精巧的设计，可能使软件系统非常丝滑
软件世界的美
> chaTGPT 帮助思考？

Chunks head


Chunks head -> block




-
-
-
-

## 5. In-Depth: Key Components in Prometheus Storage

Memory-Mapped Files (mmap)

WAL (Write-Ahead Log)
- why? 很多列族数据库都有这个东西
  - 保证一个HA,不丢失数据



Block Compression Techniques
> block 的部分主要就是压缩了
>

TODO 几张图能解释压缩问题


Compaction and Label Standardization


Index
索引部分是为了加快读 - 很显然 - 而且是join 读 - 如何能快速的匹配到需要的sample？
Inverted Indexing in Prometheus


到这里就结束了

## 7. Best Practicves for Managing Time series data in prometheus

Label Management and Standardization
    Importance of label consistency
    Pre-processing and data merging strategies
Summary of Key Techniques
    Quick recap on block storage, compression, indexing, mmap, and compaction techniques for performance enhancement