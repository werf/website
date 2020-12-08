export default {
  created() {
    this.BECOME_SPECIALIST = "BECOME_SPECIALIST"
    this.BRING_KUBERNETES = "BRING_KUBERNETES"
    this.SOLVE_TASKS = "SOLVE_TASKS"
    this.DEMO_APP = "demo_app"
    this.OWN_APP = "own_app"
    this.config = [
      {
        code: 'basic',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.basic'),
        url: '100_basic.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'basic_build',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.basic_build'),
        url: '100_basic/10_build.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'basic_cluster',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.basic_cluster'),
        url: '100_basic/20_cluster.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'basic_deploy',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.basic_deploy'),
        url: '100_basic/30_deploy.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'basic_optimize',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.basic_optimize'),
        url: '100_basic/40_optimize.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'basic_iac',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.basic_iac'),
        url: '100_basic/50_iac.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.real_apps'),
        url: '200_real_apps.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_local',
        hidden: false,
        label: this.$t('shared.label.real_apps_local'),
        menu_label: this.$t('shared.menu_items.real_apps_local'),
        url: '200_real_apps/10_local.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_debug_iac',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.real_apps_debug_iac'),
        url: '200_real_apps/20_debug_iac.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_logging',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.real_apps_logging'),
        url: '200_real_apps/25_logging.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_assets',
        hidden: false,
        label: this.$t('shared.label.real_apps_assets'),
        menu_label: this.$t('shared.menu_items.real_apps_assets'),
        url: '200_real_apps/30_assets.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_files',
        hidden: false,
        label: this.$t('shared.label.real_apps_files'),
        menu_label: this.$t('shared.menu_items.real_apps_files'),
        url: '200_real_apps/50_files.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_email',
        hidden: false,
        label: this.$t('shared.label.real_apps_email'),
        menu_label: this.$t('shared.menu_items.real_apps_email'),
        url: '200_real_apps/60_email.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_database',
        hidden: false,
        label: this.$t('shared.label.real_apps_database'),
        menu_label: this.$t('shared.menu_items.real_apps_database'),
        url: '200_real_apps/80_database.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'real_apps_stateful',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.real_apps_stateful'),
        url: '200_real_apps/90_stateful.html',
        estimated_minutes: 0, // should be overrided
        spent_minutes: 0,     // will be changed in progress
        progresstracker_lastsent_minutes: 0,     // will be changed in progress
        completed: false      // will be changed in progress
      },
      {
        code: 'infra',
        hidden: true,
        label: '',
        menu_label: this.$t('shared.menu_items.infra'),
        url: '400_infra.html',
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

      // generic
      this._addTask('basic', {
        label: this.$t('shared.tasks.basic'),
        hidden: false,
        estimated_minutes: 5
      });
      this._addTask('basic_build', {
        label: this.$t('shared.tasks.basic_build'),
        hidden: false,
        estimated_minutes: 10
      });
      this._addTask('basic_cluster', {
        label: this.$t('shared.tasks.basic_cluster'),
        hidden: false,
        estimated_minutes: 30
      });
      this._addTask('basic_deploy', {
        label: this.$t('shared.tasks.basic_deploy'),
        hidden: false,
        estimated_minutes: 10
      });
      this._addTask('basic_optimize', {
        label: this.$t('shared.tasks.basic_optimize'),
        hidden: false,
        estimated_minutes: 20
      });
      this._addTask('basic_iac', {
        label: this.$t('shared.tasks.basic_iac'),
        hidden: false,
        estimated_minutes: 10
      });
      this._addTask('real_apps', {
        label: this.$t('shared.tasks.real_apps'),
        hidden: false,
        estimated_minutes: 5
      });

      // Локальная разработка
      if (this.chosenParts.includes('real_apps_local')) {
        this._addTask('real_apps_local', {
          label: this.$t('shared.tasks.real_apps_local'),
          estimated_minutes: 20
        });
      }

      this._addTask('real_apps_debug_iac', {
        label: this.$t('shared.tasks.real_apps_debug_iac'),
        hidden: false,
        estimated_minutes: 10
      });

      this._addTask('real_apps_logging', {
        label: this.$t('shared.tasks.real_apps_logging'),
        hidden: false,
        estimated_minutes: 10
      });

      // Сборка ассетов
      if (this.chosenParts.includes('real_apps_assets')) {
        this._addTask('real_apps_assets', {
          label: this.$t('shared.tasks.real_apps_assets'),
          estimated_minutes: 15
        });
      }

      // Работа с S3
      if (this.chosenParts.includes('real_apps_files')) {
        this._addTask('real_apps_files', {
          label: this.$t('shared.tasks.real_apps_files'),
          estimated_minutes: 10
        });
      }

      // Работа с почтой
      if (this.chosenParts.includes('real_apps_email')) {
        this._addTask('real_apps_email', {
          label: this.$t('shared.tasks.real_apps_email'),
          estimated_minutes: 10
        });
      }

      // Работа с базами данных
      if (this.chosenParts.includes('real_apps_database')) {
        this._addTask('real_apps_database', {
          label: this.$t('shared.tasks.real_apps_database'),
          estimated_minutes: 20
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

    goToPage(url) {
      window.location.href=url
    }
  }
}
