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
    .state('home',{
        url: '/',
        templateUrl: 'public/partials/home.html',
    })
}])