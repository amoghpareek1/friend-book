angular.module('app').controller('UserViewController', ['MainService', 'Notification', '$state', '$stateParams', function(MainService, Notification, $state, $stateParams) {
    var _self = this

    _self.getUser = function() {
        MainService.getUser($stateParams.userID).then(function(result) {
            if(result.Success) {
                _self.userDetails = result.Data
                console.log('called')
                console.log(_self.userDetails)
            } else {
                Notification.error(result.Data)
            }
        })
    }

    _self.getUser()

    _self.sendFriendRequest = function() {
        MainService.sendFriendRequest($stateParams.userID).then(function(result) {
            if(result.Success) {
                Notification(result.Data)
            } else {
                Notification.error(result.Data)
            }
        })
    }

    _self.sendUnfriendRequest = function() {
        MainService.sendUnfriendRequest($stateParams.userID).then(function(result) {
            if(result.Success) {
                Notification(result.Data)
            } else {
                Notification.error(result.Data)
            }
        })
    }
   
}])