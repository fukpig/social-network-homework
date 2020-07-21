# Нагрузочное тестирование
1) Написал скрипт, который с faker-ом генерирует 1000000 записей.
2) Реализовал поиск в сервисе соц. сети
3) Создал скрипт для wrk, который рандомно берет символ для name и surname
3) Запускаю тестирование 
```sh
$ docker run -v /Users/rinat/Go/src/social-network/wrk/random.lua:/random.lua --rm --net=host 1vlad/wrk2-docker -t1 -c1 -d60s -R1 --latency -L -s /random.lua http://host.docker.internal:3000/search-test &> result1.out
```
```sh
$ docker run -v /Users/rinat/Go/src/social-network/wrk/random.lua:/random.lua --rm --net=host 1vlad/wrk2-docker -t1 -c10 -d60s -R10 --latency -L -s /random.lua http://host.docker.internal:3000/search-test &> result10.out
```
```sh
$ docker run -v /Users/rinat/Go/src/social-network/wrk/random.lua:/random.lua --rm --net=host 1vlad/wrk2-docker -t1 -c100 -d60s -R100 --latency -L -s /random.lua http://host.docker.internal:3000/search-test &> result100.out
```
```sh
$ docker run -v /Users/rinat/Go/src/social-network/wrk/random.lua:/random.lua --rm --net=host 1vlad/wrk2-docker -t1 -c1000 -d60s -R1000 --latency -L -s /random.lua http://host.docker.internal:3000/search-test &> result1000.out
```
# График latency до индекса
![N|Solid](https://github.com/fukpig/social-network-homework/blob/master/wrk/before_index.png?raw=true)

## Данные по результатам тестов
### 1 одновременный запрос
```sh 
51 requests in 1.00m, 16.30MB read Requests/sec: 0.85 Transfer/sec: 278.08KB
```
### 10 одновременный запрос
```sh 
138 requests in 1.00m, 41.26MB read Socket errors: connect 0, read 0, write 0, timeout 162 Requests/sec: 2.30 Transfer/sec: 703.30KB
```

### 100 одновременный запрос
```sh 
100 requests in 1.01m, 27.44MB read Socket errors: connect 0, read 0, write 0, timeout 2900 Requests/sec: 1.65 Transfer/sec: 464.37KB
```

### 1000 одновременный запрос
```sh 
110 requests in 1.02m, 29.44MB read Socket errors: connect 0, read 0, write 3, timeout 27890 Requests/sec: 1.80 Transfer/sec: 494.06KB
```

# Проставляем индекс
```sh  
ALTER TABLE `app`.`users` ADD INDEX `idx` (`name`, `surname`, `id`);
```

# Смотрим explain

```json
 {
   "query_block":{
      "select_id":1,
      "cost_info":{
         "query_cost":"224328.80"
      },
      "ordering_operation":{
         "using_filesort":false,
         "table":{
            "table_name":"users",
            "access_type":"index",
            "possible_keys":[
               "idx"
            ],
            "key":"PRIMARY",
            "used_key_parts":[
               "id"
            ],
            "key_length":"8",
            "rows_examined_per_scan":1046839,
            "rows_produced_per_join":22514,
            "filtered":"2.15",
            "cost_info":{
               "read_cost":"219825.92",
               "eval_cost":"4502.88",
               "prefix_cost":"224328.80",
               "data_read_per_join":"31M"
            },
            "used_columns":[
               "id",
               "email",
               "name",
               "surname",
               "password",
               "sex",
               "city",
               "interests"
            ],
            "attached_condition":"((`app`.`users`.`name` like 'a%') and (`app`.`users`.`surname` like 'b%'))"
         }
      }
   }
}
```

# Запускаем тесты заново

# График latency до индекса
![N|Solid](https://github.com/fukpig/social-network-homework/blob/master/wrk/after_index.png?raw=true)


## Данные по результатам тестов после индекса
### 1 одновременный запрос
```sh 
 60 requests in 1.00m, 17.64MB read Requests/sec: 1.00
```
### 10 одновременный запрос
```sh 
598 requests in 1.00m, 177.45MB read Socket errors: connect 0, read 0, write 0, timeout 9 Requests/sec: 9.96 Transfer/sec: 2.95MB
```

### 100 одновременный запрос
```sh 
1088 requests in 1.00m, 292.98MB read Socket errors: connect 0, read 0, write 0, timeout 2132 Requests/sec: 18.06
```


