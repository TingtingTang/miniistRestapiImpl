var express = require('express');
var router = express.Router();
/*
* 定义数据获取api
*
* */

router.get('/login/', function(req, res) {
    console.log("this is test")
    res.send("测试一下API的使用")
});


module.exports = router;