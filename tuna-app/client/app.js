// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	$("#error_create").hide();
	
	$scope.queryAllPrfixes = function(){

		var prefixes = $scope.prefixes;

		appFactory.queryAllPrfixes(prefixes, function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_prefixes = array;
		});
	}

	$scope.queryPrefix = function(){

		var prefix = $scope.ipPrefix;
		console.log(prefix)

		appFactory.queryPrefix(prefix, function(data){
			$scope.query_prefix = data;

			if ($scope.query_prefix == "Could not query IPv6 prefix"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.allocatePrefix = function(){

		appFactory.allocatePrefix($scope.prefix, function(data){
			$scope.allocate_prefix = data;
			if ($scope.allocate_prefix == "Could not allocate this IPv6 prefix" || $scope.allocate_prefix.indexOf("Failed to allocate IPv6 prefix") >= 0){
				console.log()
				$("#error_create").show();
				$("#success_create").hide();
			} else{
				$("#error_create").hide();
				$("#success_create").show();
			}
		});
	}

	$scope.changeHolder = function(){
		console.log($scope.holder)
		appFactory.changeHolder($scope.holder, function(data){
			$scope.change_holder = data;
			// if ($scope.change_holder == "Error: no prefix found"){
			// 	$("#error_holder").show();
			// 	$("#success_holder").hide();
			// } else{
			// 	$("#success_holder").show();
			// 	$("#error_holder").hide();
			// }
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllPrfixes = function(data, callback){

    	var prefixes = data.startKey.replace('/', '@') + '-' + data.endKey.replace('/', '@');

    	$http.get('/get_all_prefixes/'+prefixes).success(function(output){
			callback(output)
		});
	}

	factory.queryPrefix = function(prefix, callback){

		var address = prefix.replace('/', '@')

    	$http.get('/get_prefix/'+address).success(function(output){
			callback(output)
		});
	}

	factory.allocatePrefix = function(data, callback){

		var ipPrefx = data.IPv6_prefix.replace('/', '@');
		var prefix = ipPrefx + "-" + data.AS_number + "-" + data.Assign_by + "-" + data.Assign_to + "-" + data.Advertisement + "-" + new Date().getTime();
    	$http.get('/allocate_prefix/'+prefix).success(function(output){
			callback(output)
		});
	}

	factory.changeHolder = function(data, callback){

		var holder = data.prefix.replace('/', '@') + "-" + data.assignTo + "-" + new Date().getTime();

    	$http.get('/change_holder/'+holder).success(function(output){
			callback(output)
		});
	}

	return factory;
});
