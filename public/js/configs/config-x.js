angular.module('app').config(['$locationProvider', '$stateProvider', '$urlRouterProvider', 'NotificationProvider', function($locationProvider, $stateProvider, $urlRouterProvider, NotificationProvider) {
    $locationProvider.html5Mode(true)
    
    $urlRouterProvider.otherwise('/')

    NotificationProvider.setOptions({
        replaceMessage: true,
    })

    $stateProvider
    .state('signOut', {
        onEnter: function() {
            window.location = window.location.origin + '/api/v1/sign-out'
        }
    })
    .state('user',{
        url: '/',
        templateUrl: 'public/partials/user.html',
        redirectTo: 'user.list'
    })

    .state('user.list',{
        url: 'list',
        templateUrl: 'public/partials/user-list.html',
        controller: 'UsersListController as ctrl'
    })

    .state('user.edit',{
        url: 'edit',
        templateUrl: 'public/partials/user-edit.html',
        controller: 'UserEditController as ctrl'
    })
}])