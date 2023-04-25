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

  tippy('span.tooltip-text', {
    content(reference) {
      const title = reference.getAttribute('title');
      reference.removeAttribute('title');
      return title;
    },
    allowHTML: true,
    interactive: true,
    delay: [0, 200],
  });

  tippy('[data-tooltip-content]', {
    content(reference) {
      const title = reference.getAttribute('data-tooltip-content');
      reference.removeAttribute('data-tooltip-content');
      return title;
    },
    theme: 'werf',
    sticky: true,
    hideOnClick: false,
    showOnCreate: true,
    trigger: 'manual',
    allowHTML: true,
    zIndex: 4,
    popperOptions: {
      strategy: 'fixed',
      modifiers: [
        {
          name: 'flip',
          options: {
            fallbackPlacements: [],
          },
        },
        {
          name: 'preventOverflow',
          options: {
            mainAxis: false,
          },
        },
      ],
    },
  });
})