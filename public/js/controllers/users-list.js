angular.module('app').controller('UsersListController', ['MainService', 'Notification', '$state', function(MainService, Notification, $state) {
    var _self = this

	_self.filters = {
		'SortBy': 'DESC',
        'OrderBy': 'created_at',
        'searchBy' : '',
        'limit': 20,
        'offset':0
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
            console.log(_self.filters)
            if(result.Success) {
                console.log('inside')
                _self.users = result.Data.Users
                _self.totalUsersCount = result.Data.TotalCount
                console.log(_self.users)
            } else {
                console.log(result)
                Notification.error
            }
        })
    }

    _self.getUsers()

    $(document).ready(function(){
        $('#deleteButton').click(function(){
           $('.ui.basic.modal').modal('show');    
        });
    });

    _self.deleteUserModalVisible = false
    _self.requestInProgress = false
    
    _self.toggleDeleteUserModal = function() {
        _self.deleteUserModalVisible = !_self.deleteUserModalVisible
        console.log(_self.deleteUserModalVisible)
	}

	_self.deleteUser = function() {
		_self.requestInProgress = true
		MainService.deleteUser().then(function (result) {
			if(result.Success) {
				Notification(result.Data)
				$state.go('signOut')
			} else {
				Notification.error(result.Data)
			}
			_self.requestInProgress = false
		})
    }
}])