document.addEventListener('DOMContentLoaded', () => {
  tippy('#installation__release-channels span', {
    content(reference) {
      const title = reference.getAttribute('title');
      reference.removeAttribute('title');
      return title;
    },
    allowHTML: true,
    interactive: true,
    delay: [0, 200],
  });

  tippy('#installation__version span', {
    content(reference) {
      const title = reference.getAttribute('title');
      reference.removeAttribute('title');
      return title;
    },
    allowHTML: true,
    interactive: true,
    delay: [0, 200],
  });
})