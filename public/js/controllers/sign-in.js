angular.module('app').controller('SignInController', ['MainService', 'Notification', function(MainService, Notification) {
    var _self = this

    _self.requestInProgress = false

    _self.formData = {
        'Email' : '',
        'Password' : ''
    }

    _self.signIn = function(){
        _self.requestInProgress = true
        MainService.signIn(_self.formData).then(function(result){
            if(result.Success){
                Notification(result.Data)
                window.location.reload()
            }else{
                Notification.error(result.Data)
            }
            _self.requestInProgress = false
        })
    }
}])