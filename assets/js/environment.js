document.addEventListener('DOMContentLoaded', () => {
  const environment = document.getElementById('preparing-the-environment');
  const environmentRu = document.getElementById('подготовка-окружения');
  const button = document.querySelector('.sidebar-button--environment');
  const lang = document.documentElement.lang;
  const yOffset = -90;

  if (lang === 'en') {
    if (environment) {
      button.setAttribute('href', '#preparing-the-environment');
      button.removeAttribute('target');

      button.addEventListener('click', (e) => {
        e.preventDefault();

        const y = environment.getBoundingClientRect().top + window.pageYOffset + yOffset;
        window.scrollTo({top: y, behavior: 'smooth'});
      })
    }
  } else {
    if (environmentRu) {
      button.setAttribute('href', '#подготовка-окружения');
      button.removeAttribute('target');
  
      button.addEventListener('click', (e) => {
        e.preventDefault();
  
        const yRu = environmentRu.getBoundingClientRect().top + window.pageYOffset + yOffset;
        window.scrollTo({top: yRu, behavior: 'smooth'});
      })
    }
  }
})