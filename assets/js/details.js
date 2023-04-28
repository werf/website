$( document ).ready(function() {
  // Toggling details in tabs.
  $(document).on('click tap', '.details__summary', (e) => {
    e.preventDefault();
    $(e.target).closest('.details').toggleClass('active');
  });

});
