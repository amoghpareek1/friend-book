angular.module('app').controller('SignUpController', ['MainService', 'Notification', '$state', function(MainService, Notification, $state) {
    var _self = this

    _self.formData = {
        'Name'    : '',
        'Phone': '',
        'Password': '',
        'Email': ''
    }

    _self.submit = function() {
        _self.requestInProgress = true
        MainService.signUp(_self.formData).then(function(result) {
            console.log(result)
            if(result.Success) {
                Notification(result.Data)
                $state.go('signIn')
            } else {
                Notification.error(result.Data)
            }
            _self.requestInProgress = false
        })
    }
}])