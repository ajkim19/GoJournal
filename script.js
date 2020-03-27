$("form").on('submit', function(e) {
    e.preventDefault();
      
    $(".description").val($(".description").val().replace(/\r\n|\r|\n/g,"<br />"));
  });

document.getElementById("edit-entry-date").style.display = "none";