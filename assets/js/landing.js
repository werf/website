window.onload = () => {
  const triggerElements = document.querySelectorAll('[data-sm-trigger]');
  const captainElement = document.querySelector('[data-sm-captain]');

  const observer = new IntersectionObserver((entries, observer) => {
      entries.forEach(entry => {
          if (entry.isIntersecting) {
              captainElement.setAttribute('data-sm-captain', entry.target.getAttribute('data-sm-trigger'))
          }
      })
  }, { threshold: 0.9 })

  triggerElements.forEach(el => observer.observe(el))
}