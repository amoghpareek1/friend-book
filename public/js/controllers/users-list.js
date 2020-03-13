angular.module('app').controller('UsersListController', ['MainService', 'Notification', '$state', function(MainService, Notification, $state) {
    var _self = this

	_self.filters = {
		'SortBy': 'DESC',
		'OrderBy': 'created_at'
    }

    _self.toggleOrderBy = function(orderBy) {
        if(_self.filters.OrderBy === orderBy) {
            if(_self.filters.SortBy === 'DESC') {
                _self.filters.SortBy = 'ASC'
            } else {
                _self.filters.SortBy = 'DESC'
            }
        } else {
            _self.filters.SortBy = 'DESC'
        }
        _self.filters.OrderBy = orderBy
        
        _self.getUsers()
    }

    _self.getUsers = function() {
        MainService.getUsers(_self.filters).then(function(result) {
            if(result.Success) {
                console.log('inside')
                _self.users = result.Data
                console.log(_self.users)
            } else {
                console.log(result)
                Notification.error
            }
        })
    }

    _self.getUsers()

}])