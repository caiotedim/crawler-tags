$(document).ready(function() {
  $('#topfollowers').DataTable( {
      columnDefs: [
          {
              targets: [ 0, 1],
              className: 'mdl-data-table__cell--non-numeric'
          }
      ]
  } );
} );

var crawlertagsApp = angular.module('crawlertagsApp', []);
crawlertagsApp.controller('UsersList', function ($scope, $http) {
 //$http.get('/api/topfollowers').success(function(data) {
  $http({
    method  : 'GET',
    url     : '/api/topfollowers',
    data    : $scope.users,
    headers : { 'Content-Type': 'application/json' }
  }).then(function(response) {
    //
  }).catch(function(response) {
    //
  });
  $scope.orderProp = 'followers_count';
});
