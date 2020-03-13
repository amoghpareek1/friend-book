angular.module('app').service('MainService', ['$http', function($http) {
	var _self = this

	_self.signUp = function(formData) {
		return $http.post('/api/v1/sign-up', formData).then(function(response) {
			return response.data
		})
	}

	_self.signIn = function(formData) {
		return $http.post('/api/v1/sign-in',formData).then(function(response){
			return response.data
		})
	}
}])