var ctrl = {

  // Fetch another imamge 
  nextImage: function() {
      var xmlHttp = new XMLHttpRequest();
      xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
          var resp = xmlHttp.responseText;
          var holder = document.getElementById("image_holder");
          console.log("resp",resp);
          holder.setAttribute("src", "/static/images/" + JSON.parse(resp).ImageUrl);
        }
      }
      xmlHttp.open("GET", "/next_image/", true);
      xmlHttp.send(null);
  },
  
  // Handle like 
  handleLike : function() {
    this.nextImage();
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "/preference_event/", true);
    xmlHttp.send(JSON.stringify({"Liked": "true"}));
  },
  
  // Handle dislike
  handleDislike : function() {
    this.nextImage();
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "/preference_event/", true);
    xmlHttp.send(JSON.stringify({"Liked": "false"}));
  }

};
