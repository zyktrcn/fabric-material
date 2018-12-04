

var prefix = require('./controller.js');

module.exports = function(app){

  app.get('/get_prefix/:address', function(req, res){
    prefix.get_prefix(req, res);
  });
  app.get('/allocate_prefix/:prefix', function(req, res){
    prefix.allocate_prefix(req, res);
  });
  app.get('/get_all_prefixes/:prefixes', function(req, res){
    prefix.get_all_prefixes(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    prefix.change_holder(req, res);
  });
}