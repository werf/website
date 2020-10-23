module.exports = {
  'server': {
    baseDir: '_site',
    port: 3000,
    routes: {
      "/applications_guide_ru": ""
    }
  },
  'files': '_site',
  'serveStatic': [
    {
      route: "/applications_guide_ru",
      dir: '_site'
    }
  ],
  'serveStaticOptions': {
    'extensions': ['html']
  }
}
