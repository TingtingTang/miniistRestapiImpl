var express = require('express');

var app = express();
// set up handlebars view engine
var handlebars = require('express-handlebars').create({ defaultLayout:'main' });
app.engine('handlebars', handlebars.engine);
app.set('view engine', 'handlebars');

app.set('port', process.env.PORT || 3001);

//before declaring any route(s) - establish static directory
app.use(express.static(__dirname + '/app'));

app.use(function(req, res, next){
    res.locals.showTests = app.get('env') !== 'production' && req.query.test === '1';
    next();
});

//routes
var approutes = require('./approutes.js')(app);
var api = require('./handlers/api');
app.use('/api', api);


//allow custom header and CORS (2016-02-28)
app.all('*', function (req, res, next) {
    res.header('Access-Control-Allow-Origin', '*');

    if (req.method == 'OPTIONS') {
        res.send(200);          
    } else {
        next();
    }
});

app.listen(app.get('port'), function(){
    console.log( 'Express started on http://localhost:' +
    app.get('port') + '; press Ctrl-C to terminate.' );
});

