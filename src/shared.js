export default {
  created() {
    this.BECOME_SPECIALIST = "BECOME_SPECIALIST"
    this.BRING_KUBERNETES = "BRING_KUBERNETES"
    this.SOLVE_TASKS = "SOLVE_TASKS"
    this.DEMO_APP = "demo_app"
    this.OWN_APP = "own_app"
    this.config = [
      {
        code: 'hasKubeCluster',
        label: 'У меня есть готовый Kubernetes-кластер, связанный с репозиторием',
        menu_label: 'Подготовка к работе',
        url: '010_preparing.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'build',
        label: 'Сборка образа приложения',
        menu_label: 'Сборка',
        hidden: true,
        url: '020_basic/10_build.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'IaC',
        label: 'Описание архитектуры приложения как код и деплой',
        menu_label: 'Конфигурирование инфраструктуры в виде кода',
        hidden: true,
        url: '020_basic/20_iac.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'ci',
        label: 'Построение CI-процесса',
        menu_label: 'Построение CI-процесса',
        hidden: true,
        url: '020_basic/30_ci.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'usingDependencyManagers',
        menu_label: 'Подключаем зависимости',
        label: 'В моих приложениях обычно используются менеджеры зависимостей',
        url: '030_dependencies.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'hasAssets',
        label: 'В моих приложениях обычно есть ассеты (css, js) и их сборка',
        menu_label: 'Генерируем и раздаём ассеты',
        url: '040_assets.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'hasUserGeneratedFiles',
        label: 'Мои приложения обычно хранят файлы на диск',
        menu_label: 'Работа с файлами',
        url: '050_files.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'hasEmail',
        label: 'Моим приложения нужна отправка данных по e-mail',
        menu_label: 'Работа с электронной почтой',
        url: '060_email.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'hasInMemoryDB',
        label: 'Я хочу использовать in-memory базу данных',
        menu_label: 'Подключаем Redis',
        url: '070_redis.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'hasAnyDB',
        label: 'Я хочу использовать базу данных',
        menu_label: 'Подключаем базу данных',
        url: '080_database.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'unitTesting',
        label: 'Я хочу использовать юнит-тесты или линтеры',
        menu_label: 'Юнит-тесты и Линтеры',
        url: '090_unittesting.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'multipleApps',
        label: 'Я обычно храню несколько приложений в одном репозитории',
        menu_label: 'Несколько приложений в одном репозитории',
        url: '110_multipleapps.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'dynamicEnvs',
        label: 'Хочу иметь возможность выкатывать каждую фичу на отдельный контур',
        menu_label: 'Динамические окружения',
        url: '120_dynamicenvs.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
    ];
  },

  props: {
    framework: String,
    ci: String,
    guideCode: String,
    baseUrl: String,
  },

  data() {
    return {
      skippedForm: localStorage.skippedForm !== 'false',
      wish: localStorage.wish?localStorage.wish:"",             // "become_specialist" | "bring_kubernetes" | "solve_tasks" // @see created()
      usingApp: localStorage.usingApp?localStorage.usingApp:"", // "demo_app" | "own_app" // @see created()
      learningProgress: localStorage.learningProgress?JSON.parse(localStorage.learningProgress):[], // array of elements from this.config, overrided in recalcForm()

      completedStepsCalculated: 0,
      allCompleteCalculated: false,
      nextTaskCalculated: {url: "", menu_label: ""},
      currentPage: this._getCurUrl(),
    }
  },

  mounted() {
    this.recalcTasks()
  },

  methods: {
    _saveLocalStorage() {
      localStorage.skippedForm = this.skippedForm
      localStorage.wish = this.wish
      localStorage.usingApp = this.usingApp
      localStorage.learningProgress = JSON.stringify(this.learningProgress)
    },

    _getCurUrl() {
      let pattern = new RegExp('\/'+this.guideCode+'\/(.+)');
      let match = pattern.exec(window.location.pathname);
      if (match)
        return match[1];
      return "undefined";
    },

    /* ------------------------------------------------------------------------------------------------------------ */

    _addTask(code, override = {}) {
      // anti-duplication
      // console.log(this.learningProgress);
      // if (this.learningProgress && this.learningProgress.find(el => el.code===code))
      //   return;
      // find and add element
      let el = this.config.find(el => el.code === code)
      el = {...el, ...override};
      this.learningProgress.push(el);
      // recalc minutes-hours-days
      this.estimatedMinutesCalculated += el.estimated_minutes;
      this.estimatedHoursCalculated = Math.ceil(this.estimatedMinutesCalculated / 60)
      this.estimatedDaysCalculated = Math.ceil(this.estimatedHoursCalculated / 8)
      // Save to localStorage
      this._saveLocalStorage()
    },

    markCurrentPageComplete(event){
      for (let i in this.learningProgress) {
        if (this.learningProgress[i].url === this.currentPage) {
          this.learningProgress[i].completed = true;
          // TODO: тут надо стопать таймер насовсем
          this._saveLocalStorage()
        }
      }
    },

    _getCurrentPageTimer() {
      if (! this.learningProgress)
        return 0;
      let $this = this;
      let current_task = this.learningProgress.find(task => task.url===$this.currentPage)
      if (! current_task)
        return 0;
      return current_task.spent_minutes
    },

    _setCurrentPageTimer(value) {
      if (! this.learningProgress)
        return false;
      let $this = this;
      let current_task = this.learningProgress.find(task => task.url===$this.currentPage)
      if (! current_task)
        return false;
      current_task.spent_minutes = value
      this._saveLocalStorage()
    },

    _getCurrentPageEstimate() {
      if (! this.learningProgress)
        return 0;
      let $this = this;
      let current_task = this.learningProgress.find(task => task.url===$this.currentPage)
      if (! current_task)
        return 0;
      return current_task.estimated_minutes
    },

    _getCurrentPageCode() {
      if (! this.learningProgress)
        return false;
      let $this = this;
      let current_task = this.learningProgress.find(task => task.url===$this.currentPage)
      if (! current_task)
        return false;
      return current_task.code;
    },

    _getCurrentPageTrackerMinutes() {
      if (! this.learningProgress)
        return false;
      let $this = this;
      let current_task = this.learningProgress.find(task => task.url===$this.currentPage)
      if (! current_task)
        return false;
      return current_task.progresstracker_lastsent_minutes;
    },

    _setCurrentPageTrackerMinutes(value) {
      if (! this.learningProgress)
        return false;
      let $this = this;
      let current_task = this.learningProgress.find(task => task.url===$this.currentPage)
      if (! current_task)
        return false;
      current_task.progresstracker_lastsent_minutes = value;
    },

    /* ------------------------------------------------------------------------------------------------------------ */

    _clearLearningProcess() {
      this.learningProgress = []
    },

    recalcForm() {
      this.estimatedMinutesCalculated = 0;

      this.learningProgress = []
      // hasKubeCluster
      if (this.chosenParts.includes('hasKubeCluster')) {
        this._addTask('hasKubeCluster', {
          label: 'понимание требований к инфраструктуре',
          estimated_minutes: 30
        });
      } else {
        this._addTask('hasKubeCluster', {
          label: 'настройка инфраструктуры pet-проекта',
          estimated_minutes: 120
        });
      }
      // build, IaC, ci - should always be
      this._addTask('build', {
        hidden: false,
        estimated_minutes: 60
      });
      this._addTask('IaC', {
        hidden: false,
        estimated_minutes: 90
      });
      this._addTask('ci', {
        hidden: false,
        estimated_minutes: 30
      });
      // usingDependencyManagers
      if (this.chosenParts.includes('usingDependencyManagers')) {
        this._addTask('usingDependencyManagers', {
          label: 'описание инфраструктуры для обновления зависимостей',
          estimated_minutes: 45
        });
      }
      // hasAssets
      if (this.chosenParts.includes('hasAssets')) {
        this._addTask('hasAssets', {
          label: 'описание инфраструктуры для сборки ассетов',
          estimated_minutes: 30
        });
      }
      // hasUserGeneratedFiles
      if (this.chosenParts.includes('hasUserGeneratedFiles')) {
        this._addTask('hasUserGeneratedFiles', {
          label: 'описание инфраструктуры для хранения файлов',
          estimated_minutes: 15
        });
      }
      // hasEmail
      if (this.chosenParts.includes('hasEmail')) {
        this._addTask('hasEmail', {
          label: 'понимание, как работать с e-mail',
          estimated_minutes: 15
        });
      }
      // hasInMemoryDB
      if (this.chosenParts.includes('hasInMemoryDB')) {
        this._addTask('hasInMemoryDB', {
          label: 'понимание, как работать с in-memory базой данных',
          estimated_minutes: 60
        });
      }
      // hasAnyDB
      if (this.chosenParts.includes('hasAnyDB')) {
        this._addTask('hasInMemoryDB', {
          label: 'понимание, как работать с in-memory базой данных',
          estimated_minutes: 60
        });
        this._addTask('hasAnyDB', {
          label: 'понимание, как работать с базами данных',
          estimated_minutes: 90
        });
      }
      // unitTesting
      if (this.chosenParts.includes('unitTesting')) {
        this._addTask('unitTesting', {
          label: 'запуск и реализация юнит-тестов и линтеров',
          estimated_minutes: 30
        });
      }
      // multipleApps
      if (this.chosenParts.includes('multipleApps')) {
        this._addTask('multipleApps', {
          label: 'организация монорепозитория',
          estimated_minutes: 45
        });
      }
      // dynamicEnvs
      if (this.chosenParts.includes('dynamicEnvs')) {
        this._addTask('dynamicEnvs', {
          label: 'организация инфраструктуры для динамических окружений',
          estimated_minutes: 15
        });
      }
    },

    recalcTasks() {
      this.completedStepsCalculated = this.learningProgress.filter(step => step.completed).length;
      if (this.completedStepsCalculated===this.learningProgress.length) {
        this.allCompleteCalculated = true;
        this.nextTaskCalculated = null;
      } else {
        this.allCompleteCalculated = false;
        let $this = this;
        this.nextTaskCalculated = this.learningProgress.find(task => !task.completed && (task.url!==$this.currentPage))
      }
    },
  }
}
