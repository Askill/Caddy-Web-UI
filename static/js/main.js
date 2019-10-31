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

function loadSites(){
    getJSON('http://localhost:8000/api/Sites',
        function(err, data) {
            if (err !== null) {
                alert('Something went wrong: ' + err);
            } else {
                for (index = 0; index < Object.keys(data).length; index++) { 
                    
                    renderSite(data[index]);
                } 
            }
        }
    );
}

function renderSite(site){
    site = ``

    console.log(site); 
}