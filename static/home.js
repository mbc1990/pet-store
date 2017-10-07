var ctrl = {

  // Fetch another imamge 
  nextImage: function() {
      var xmlHttp = new XMLHttpRequest();
      xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
          var resp = xmlHttp.responseText;
          var holder = document.getElementById("image_holder");
          console.log("resp",resp);
          holder.setAttribute("src", JSON.parse(resp).ImageUrl);
        }
      }
      xmlHttp.open("GET", "/next_image/", true);
      xmlHttp.send(null);
  },
  
  // Handle like 
  handleLike : function() {
    console.log("Like clicked");
    this.nextImage();
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "/preference_event/", true);
    xmlHttp.send(JSON.stringify({"liked": true}));
  },
  
  // Handle dislike
  handleDislike : function() {
    console.log("Dislike clicked");
    this.nextImage();
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "/preference_event/", true);
    xmlHttp.send(JSON.stringify({"liked": false}));
  }

};
