$(document).ready(function () {
  const $uploadContainer = $('#upload-container');
  const $fileInput = $('#file-input');
  const $uploadBtn = $('#upload-btn');

  // Button click triggers file input
  $uploadBtn.on('click', function () {
    $fileInput.click();
  });

  // Drag over effect
  $uploadContainer.on('dragover', function (e) {
    e.preventDefault();
    e.stopPropagation();
    $uploadContainer.addClass('dragover');
  });

  $uploadContainer.on('dragleave', function (e) {
    e.preventDefault();
    e.stopPropagation();
    $uploadContainer.removeClass('dragover');
  });

  // Handle file drop
  $uploadContainer.on('drop', function (e) {
    e.preventDefault();
    e.stopPropagation();
    $uploadContainer.removeClass('dragover');

    const file = e.originalEvent.dataTransfer.files[0];
    if (file) {
      uploadFile(file);
    }
  });

  // Handle file input selection
  $fileInput.on('change', function () {
    const file = this.files[0];
    if (file) {
      uploadFile(file);
    }
  });

  function uploadFile(file) {
    const formData = new FormData();
    formData.append("file", file); 

    fetch('/api/process', {
      method: 'POST',
      credentials: "include",
      body: formData,
    })
      .then(response => response.text())
      .then(result => console.log(result))
      .catch(error => console.error('Error:', error));
  }
});
