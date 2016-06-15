var express = require('express');
var router = express.Router();
/*
* 定义数据获取api
*
* */

//router.get('/login/', function(req, res) {
//    console.log("this is test")
//    res.send("测试一下API的使用")
//});
var req_label = {
    "action": "get_labelStatis"
};
router.get('/statistics/summary', function (req, res) {
    var response = {
        "machine": 12,
        "test_suit": 70,
        "workload": 26,
        "tasks": 60,
        "ret": 0,
        "info": "OK"
    };
    if(req_label.action == "get_labelStatis"){
        res = response;
    }
});
//var server = express.listen(8081, function(){
//    var host = server.address().address;
//    var port = server.address().port;
//});

$(function () {
    function getLabels() {
        $.ajax({
            url: "api/statistics/summary",
            type: "get",
            data: req_label,
            dataType: "json",
            success: function (data) {
                console.info(data);
            }
        })
    }
    getLabels();
});

module.exports = router;