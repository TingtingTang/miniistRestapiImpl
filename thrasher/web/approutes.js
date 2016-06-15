var main = require('./handlers/main.js');

module.exports = function(app){
    app.get('/', main.index);
    app.get('/index', main.index);
    app.get('/home', main.index);
};
