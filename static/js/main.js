var getJSON = function(url, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', url, true);
    xhr.responseType = 'json';
    xhr.onload = function() {
      var status = xhr.status;
      if (status === 200) {
        callback(null, xhr.response);
      } else {
        callback(status, xhr.response);
      }
    };
    xhr.send();
};
var postJSON = function(url, data, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.responseType = 'json';
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
      var status = xhr.status;
      if (status === 200) {
        callback(null, xhr.response);
      } else {
        callback(status, xhr.response);
      }
    };
    xhr.send(data);
};
var deleteJSON = function(url, data, callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('DELETE', url + "/" + data, true);
    xhr.responseType = 'json';
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
      var status = xhr.status;
      if (status === 200) {
        callback(null, xhr.response);
      } else {
        callback(status, xhr.response);
      }
    };
    xhr.send();
};

function loadSites(){
    getJSON('http://localhost:8000/api/Sites',
        function(err, data) {
            data = Object.values(data);
            if (err !== null) {
                alert('Something went wrong: ' + err);
            } else {
                $(".rendered").remove();
                data.forEach (function(site) { 
                    renderSite(site);
                })
            }
        }
    );
}

function renderSite(site){
    siteString = `
        <div  class="card  float-sm-left border-success lg-3 md-4 sm-2 rendered" id="${site.id}" >
            <div class="card-header"><a href="https://${site.domain}">${site.domain}</a> </div>
            <div class="card-body">
                <h4 class="card-title">${site.title}</h4>

                <h6 class="card-text">${site.source} -> ${site.target}</p>
                <h6 class="card-text">tls: ${site.email}</p>
                <p class="card-text">${site.description}</p>
            </div>
            <button type="button" class="btn btn-outline-warning float-right" onclick="deleteSite('${site.id}')">Delete</button>
        </div>
    `
    document.getElementById("main").innerHTML += siteString;
}

function newSite(){
    let json = {};
    json["domain"] = document.getElementById('domain').value;
    json["title"] = document.getElementById('title').value;
    json["description"] = document.getElementById('description').value;
    json["source"] = document.getElementById('source').value;
    json["target"] = document.getElementById('target').value;
    json["email"] = document.getElementById('email').value;
    console.log(json);
    postJSON("http://localhost:8000/api/Sites", JSON.stringify(json), loadSites);
}

function deleteSite(id){
    deleteJSON("http://localhost:8000/api/Sites", id, loadSites)
}