angular.module('app').controller('UserEditController', ['MainService', 'Notification', '$state', function(MainService, Notification, $state) {
    var _self = this

    _self.requestInProgress = false

   _self.userDetails = []


    MainService.getMe().then(function(result){
        if(result.Success){
            _self.userDetails = result.Data
        } else {
            Notification.error(result.Data)
        }
    })

    _self.save = function() {
        _self.requestInProgress = true
        MainService.putUserDetails(_self.userDetails).then(function(result){
            if(result.Success){
                Notification(result.Data)
                $state.go('user.list')
            } else {
                Notification.error(result.Data)
            }
            _self.requestInProgress = false
        })
    }
}])