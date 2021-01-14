$( document ).ready(function() {
  $('.details__summary').on('click tap', function() {
    $(this).closest('.details').toggleClass('active');
  });

  $('.open-content_button').on('click tap', function() {
    $(this).closest('.open-content_container').toggleClass('active');
  });
});