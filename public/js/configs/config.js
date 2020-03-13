angular.module('app').config(['$locationProvider', '$stateProvider', '$urlRouterProvider', 'NotificationProvider', function($locationProvider, $stateProvider, $urlRouterProvider, NotificationProvider) {
    $locationProvider.html5Mode(true)
    
    $urlRouterProvider.otherwise('/')

    NotificationProvider.setOptions({
        replaceMessage: true
    })

    $stateProvider
    .state('signIn', {
        url: '/',
        templateUrl: 'public/partials/sign-in.html',
        controller: 'SignInController as ctrl'
    })

    .state('signUp', {
        url: '/sign-up',
        templateUrl: 'public/partials/sign-up.html',
        controller: 'SignUpController as ctrl'
    })
}])