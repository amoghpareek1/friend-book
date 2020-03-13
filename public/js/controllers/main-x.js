angular.module('app').controller('MainXController', ['MainService', function (MainService) {
	var _self = this

	_self.me = {}
	MainService.getMe().then(function (result) {
		if (result.Success) {
			_self.me = result.Data
			console.log(_self.me)
		}
	})
}])