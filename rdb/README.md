# rdb
suport redis RDB format, in order to support migrate (restore) <-> redis.
1. support redis RDB version  6 <b>dump encode</b>.
2. support redis RDB version  1 <= version <= 10(Redis 7.0) <b>parse decode</b>.
 
# reference
* [RDB_Version_History](https://github.com/sripathikrishnan/redis-rdb-tools/blob/master/docs/RDB_Version_History.textile)
* [<b><u>cupcake/rdb</u></b>](https://github.com/cupcake/rdb)
* [tair-opensource/RedisShake](https://github.com/tair-opensource/RedisShake)
* [HDT3213/rdb](https://github.com/HDT3213/rdb)