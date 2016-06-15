var influxdb = require('influx');


// 建立influxdb的连接
var client = influxdb({
    hosts: ['http://onekoko.com'],
    username: 'root',
    password: 'root',
    database: 'collectd'
});

console.log(client.getHostsAvailable());

// 获取可用的数据库列表
client.getDatabaseNames( function(err,arrayDatabaseNames){
    console.log(arrayDatabaseNames);
});

// 查询influxdb
var query = 'select value from entropy_value limit 10';
client.query(['collectd'],query, function(err, results) {
    if (err){
        console.log(err)
    }
    console.log(results);
});