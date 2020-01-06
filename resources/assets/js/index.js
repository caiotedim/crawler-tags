function getAPI() {
  getTopFollowers();
  getPostsSummarized();
  getTotalPostsLang();
}

function getTopFollowers() {
  var http = new XMLHttpRequest();
  http.onreadystatechange = function() {
    if ( http.readyState == 4 ) {
      if ( http.status != 200 ) {
        // error
        console.log("Error")
      } else {
        var content = http.response;
        var obj = jQuery.parseJSON(content);
        var tableHtml = "";
        for (var i =0 ; i < obj.length; i++) {
          tableHtml += "<tr><td>"+obj[i].user+"</td><td>"+obj[i].followers_count+"</td></tr>";
        }
        document.getElementById("topfollowers").innerHTML = tableHtml;
      }
    }
  }
  http.open("GET", "/api/topfollowers", true);
  http.setRequestHeader("Content-Type", "application/json");
  http.setRequestHeader("Accept", "application/json");
  http.send();
}

function getPostsSummarized() {
  var http = new XMLHttpRequest();
  http.onreadystatechange = function() {
    if ( http.readyState == 4 ) {
      if ( http.status != 200 ) {
        // error
        console.log("Error")
      } else {
        var content = http.response;
        var obj = jQuery.parseJSON(content);
        var tableHtml = "";
        for (var i =0 ; i < obj.length; i++) {
          tableHtml += "<tr><td>"+obj[i].Hour+"</td><td>"+obj[i].Total+"</td></tr>";
        }
        document.getElementById("postsummarized").innerHTML = tableHtml;
      }
    }
  }
  http.open("GET", "/api/postsummarized", true);
  http.setRequestHeader("Content-Type", "application/json");
  http.setRequestHeader("Accept", "application/json");
  http.send();
}

function getTotalPostsLang() {
  var http = new XMLHttpRequest();
  http.onreadystatechange = function() {
    if ( http.readyState == 4 ) {
      if ( http.status != 200 ) {
        // error
        console.log("Error")
      } else {
        var content = http.response;
        var obj = jQuery.parseJSON(content);
        var tableHtml = "";
        for (var i =0 ; i < obj.length; i++) {
          tableHtml += "<tr><td>"+obj[i].Hashtag+"</td>";
          for (var v = 0 ; v < obj[i].LangCount.length; v++) {
            tableHtml += "<td>" + obj[i].LangCount[v].Lang + " = " + obj[i].LangCount[v].Total + "</td>";
          }
          tableHtml += "</tr>";
        }
        document.getElementById("totalpostslang").innerHTML = tableHtml;
      }
    }
  }
  http.open("GET", "/api/totalpostslang", true);
  http.setRequestHeader("Content-Type", "application/json");
  http.setRequestHeader("Accept", "application/json");
  http.send();
}