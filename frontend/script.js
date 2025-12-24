$(document).ready(function() {
  const $uploadContainer = $('#upload-container');
  const $fileInput = $('#file-input');
  const $uploadBtn = $('#upload-btn');

  // Button click triggers file input
  $uploadBtn.on('click', function() {
    $fileInput.click();
  });

  // Drag over effect
  $uploadContainer.on('dragover', function(e) {
    e.preventDefault();
    e.stopPropagation();
    $uploadContainer.addClass('dragover');
  });

  $uploadContainer.on('dragleave', function(e) {
    e.preventDefault();
    e.stopPropagation();
    $uploadContainer.removeClass('dragover');
  });

  // Handle file drop
  $uploadContainer.on('drop', function(e) {
    e.preventDefault();
    e.stopPropagation();
    $uploadContainer.removeClass('dragover');
    const files = e.originalEvent.dataTransfer.files;
    handleFiles(files);
  });

  // Handle file input selection
  $fileInput.on('change', function() {
    const files = this.files;
    handleFiles(files);
  });

  function handleFiles(files) {
    if (!files || files.length === 0) return;
    
    const formData = new FormData();
    
    for (let i = 0; i < files.length; i++) {
        console.log(`${files[i]}`);
        formData.append(`file${i}`, files[i]);
    }

    //fetch('https://rozetka.hopto.org:20080/process', {
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

