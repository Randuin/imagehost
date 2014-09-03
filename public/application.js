$(function() {
  var image_container = $('#image-container');

  $('html').bind('paste', function(e) {
    var matchType = /image.*/
    var clipboardData = event.clipboardData;

    if (event.clipboardData.types[0] == 'Files') {
      for (i = 0; i < (event.clipboardData.items).length; i++) {
        var dataTransferItem = event.clipboardData.items[i];
        var file = dataTransferItem.getAsFile();

        var reader = new FileReader();
        reader.onload = function(evt) {
          var img = new Image();
          img.src = evt.target.result;

          image_container.attr('src', img.src);

          uploadFile(file);
        };
        reader.readAsDataURL(file);
      }
    }
  });

  function uploadFile(file) {
    var formData = new FormData();
    formData.append('file', file);

    $.ajax({
      type: "POST",
      url: "/upload",
      enctype: "multipart/form-data",
      data: formData,
      processData: false,
      contentType: false,
      success: function(data) {
        updateURL(data); }
    });
  }

  function updateURL(url) {
    var final_url = url;
    var input_field = $('#image-url');
    input_field.val(final_url);
    input_field.focus().select();
  }
});

