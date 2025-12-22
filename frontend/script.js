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
    for (let i = 0; i < files.length; i++) {
      alert(`File uploaded: ${files[i].name}`);
    }
  }
});

