angular.module('app').controller('MainXController', ['MainService', '$state', function (MainService, $state) {
	var _self = this

	_self.me = {}
	MainService.getMe().then(function (result) {
		if (result.Success) {
			_self.me = result.Data
			console.log(_self.me)
		}
	})
}])