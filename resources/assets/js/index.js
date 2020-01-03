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

function getAPI() {
  getTopFollowers();
}

function getTopFollowers() {
  var http = new XMLHttpRequest();
  http.onreadystatechange = function() {
    if ( http.readyState == 4 ) {
      if ( http.status != 200 ) {
        // error
      } else {
        var content = http.response;
        var obj = jQuery.parseJSON(content);
        var tableHtml = "";
        for (var i =0 ; i<= obj.length; i++) {
          tableHtml += "<tr><td>"+obj[i].user+"</td><td>"+obj[i].followers_count+"</td></tr>";
        }
        document.getElementById("topfollowers").innerHTML = tableHtml;
      }
    }
    http.open("GET", "/api/topfollowers", true);
    http.setRequestHeader("Content-Type", "application/json");
    http.setRequestHeader("Accept", "application/json");
    http.send();
  }
}
