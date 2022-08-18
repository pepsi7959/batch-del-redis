# batch-del-redis
Delete Multiple Key On REDIS Sentinel


## Build
```bash
make
```

## Run
```bash
./batch-del-redis del_redis_key.txt 192.100.200.10:26379,192.100.201.10:26379,192.100.202.10:26379
```

## Example del_redis_key.txt file
> list of your redis key that need to delelete
```txt
key:1
key:2
```
