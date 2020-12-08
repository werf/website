module.exports = {
  'server': {
    baseDir: '_site',
    port: 3000,
    routes: {
      "/applications_guide_ru": "",
      "/applications_guide_en": ""
    }
  },
  'files': '_site',
  'serveStatic': [
    {
      route: "/applications_guide_ru",
      dir: 'applications_guide_ru/_site'
    },
    {
      route: "/applications_guide_en",
      dir: 'applications_guide_en/_site'
    },
  ],
  'serveStaticOptions': {
    'extensions': ['html']
  }
}
