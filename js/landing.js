$(document).ready(function() {
  var containerElement = document.querySelector('[data-sm-container]');
  var containerMonitor = scrollMonitor.createContainer(containerElement);
  var triggerElements = document.querySelectorAll('[data-sm-trigger]');
  var captainElement = document.querySelector('[data-sm-captain]');

  triggerElements.forEach(function(trigger) {
    var m = containerMonitor.create(trigger, -300);
    m.enterViewport(function() {
      captainElement.setAttribute('data-sm-captain', m.watchItem.getAttribute('data-sm-trigger'))
    });
  })
});