var ctrl = {
  
  // Image id of the currently shown image
  currentId : 0,

  // Fetch another imamge 
  nextImage: function() {
      var xmlHttp = new XMLHttpRequest();
      xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
          var resp = JSON.parse(xmlHttp.responseText);
          console.log("Setting current id to",resp.ImageId);
          // this.currentId = resp.ImageId
          updateCurrentId(resp.ImageId);
          var holder = document.getElementById("image_holder");
          holder.setAttribute("src", "/static/images/" + resp.ImageUrl);
        }
      }
      xmlHttp.open("GET", "/next_image/", true);
      xmlHttp.send(null);
  },
  
  // Handle like 
  handleLike : function() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "/preference_event/", true);
    console.log("Sending current id: ",this.currentId);
    xmlHttp.send(JSON.stringify({"Liked": "true", "ImageId": this.currentId.toString()}));
    this.nextImage();
  },
  
  // Handle dislike
  handleDislike : function() {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST", "/preference_event/", true);
    xmlHttp.send(JSON.stringify({"Liked": "false", "ImageId": this.currentId.toString()}));
    this.nextImage();
  }

};

// Update ctrl.currentId
function updateCurrentId(imgId) {
  ctrl.currentId = imgId;
}
