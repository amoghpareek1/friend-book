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

	_self.sendFriendRequest = function(userID) {
		return $http.post('/api/v1/send-friend-request/'+userID).then(function(response){
			return response.data
		})
	}

	_self.sendUnfriendRequest = function(userID) {
		return $http.delete('/api/v1/send-unfriend-request/'+userID).then(function(response){
			return response.data
		})
	}

	_self.getMe = function() {
		return $http.get('/api/v1/get-me').then(function(response) {
			return response.data
		})
	}

	_self.getUser = function(userID) {
		return $http.get('/api/v1/user/'+ userID + '/get').then(function(response) {
			return response.data
		})
	}

	_self.getUsers = function(filters) {
		return $http({
			url: '/api/v1/users/get', 
			method: "GET",
			params: filters
		}).then(function(response) {
			return response.data
		})
	}

	_self.putUserDetails = function(userDetails){
		return $http.put('/api/v1/user/update', userDetails).then(function(response){
			return response.data
		})
	}

	_self.deleteUser = function() {
		return $http.delete('/api/v1/user/delete').then(function(response) {
			return response.data
		})
	}
}])