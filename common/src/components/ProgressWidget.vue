<template>
  <div class="progresswidget">
    <div v-if="skippedForm">
      <p>{{ $t('progresswidget.you_should') }} <a :href="baseUrl+'/'+guideCode+'/000_task.html'">{{ $t('progresswidget.define_targets') }}</a> {{ $t('progresswidget.for_better_learning') }}</p>
    </div>
    <div v-if="! skippedForm">
      <div class="progresswidget__target">Цель:
        <strong v-if="wish===BECOME_SPECIALIST">{{ $t('progresswidget.target.become_specialist') }}</strong>
        <strong v-if="wish===BRING_KUBERNETES">{{ $t('progresswidget.target.bring_kubernetes') }}</strong>
        <strong v-if="wish===SOLVE_TASKS">{{ $t('progresswidget.target.solve_tasks') }}</strong>
      </div>
      <div class="progresswidget__progress">
        {{ $t('progresswidget.progress.steps_completed') }}: <span>{{completedStepsCalculated}}</span> {{ $t('of') }} {{learningProgress.length}}. {{ $t('progresswidget.progress.current_step') }}: <div id="minutes_total" :class="'progresswidget__timer '+((timer_minutes_total>=timer_estimated_minutes)?'progresswidget__timer_overtime':'')">{{ timer_minutes_total }}</div> {{ $t('of') }}
        {{ timer_estimated_minutes }} {{ $t('progresswidget.progress.minutes') }}.
        &nbsp;&nbsp;
        <div :class="'button__blue button__blue_small button__blue_inline '+((timer_minutes_total>=timer_estimated_minutes)?'button__blue_red':'')" @click="goToPage('https://t.me/werf_ru')">
          <a href="#">{{ $t('progresswidget.progress.help_needed') }}</a>
        </div>
        &nbsp;&nbsp;
        <div class="button__blue button__blue_small button__blue_inline" v-if="!expandedList" @click="expandList">
          <a href="#">{{ $t('progresswidget.progress.check_plan') }}</a>
        </div>
      </div>
      <div v-if="expandedList" class="learningprogress">
        <div class="learningprogress__checklist">
          <div v-for="el in learningProgress" :class="'learningprogress__row '+(el['completed']?'learningprogress__row_completed':'')">
            <span class="learningprogress__time">{{el['spent_minutes']}} {{ $t('of') }} {{el['estimated_minutes']}}</span>
            <span class="learningprogress__label">{{el['label']}}</span>
          </div>
        </div>
        <div class="button__blue button__blue_small button__blue_inline" @click="hideList">
          <a href="#">{{ $t('progresswidget.progress.collapse_plan') }}</a>
        </div>
        <div class="button__blue button__blue_small button__blue_inline" @click="goToPage(baseUrl+'/'+guideCode+'/000_task.html')">
          <a href="#">{{ $t('progresswidget.progress.change_plan') }}</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import sharedMixin from "../shared"
import progressTrackerMixin from "../progresstracker"
export default {
  name: 'ProgressWidget',

  mixins: [sharedMixin, progressTrackerMixin],

  data() {
    return {
      expandedList: false,

      timer_active: true,
      session_started: (new Date()).getTime(),

      timer_activities_prev: 0, // will be overrided in created() method
      timer_activities_cur: 0,
      timer_minutes_total: 0,

      timer_estimated_minutes: 0, // will be overrided in created() method
    };
  },

  created() {
    this.timer_activities_prev = this._getCurrentPageTimer()*60000;
    this.timer_estimated_minutes = this._getCurrentPageEstimate();
    setInterval(this._timerTick, 5000)
  },

  methods: {
    expandList(){
      this.expandedList = true;
    },
    hideList(){
      this.expandedList = false;
    },
    /* ---------------------------------------- timer ---------------------------------------- */
    onIdle(){
      this.timer_active = false
      this.timer_activities_prev = this.timer_activities_prev + this.timer_activities_cur
      this.timer_activities_cur = 0
    },
    /**
     * Метод, вызываемый сторонней либой
     */
    onActive(){
      this.session_started = (new Date()).getTime()
      this.timer_active = true
    },

    /**
     * Метод, обеспечивающий работы таймера. Периодически запускается и, если нужно, меняет чиселки.
     * Вызов - смотри в created()
     */
    _timerTick(){
      if (this.timer_active){
        // Расчитываем и отображаем таймер
        this.timer_activities_cur = (new Date()).getTime() - this.session_started
        this.timer_minutes_total = Math.floor((this.timer_activities_prev + this.timer_activities_cur)/60000)
        this._setCurrentPageTimer(this.timer_minutes_total)

        // проверить, не прошло ли 10 минут и если да то надо стрелять событие
        if (this._getCurrentPageTrackerMinutes() + 10 < this.timer_minutes_total) {
          this._setCurrentPageTrackerMinutes(this.timer_minutes_total);
          this._progresstracker_timer_fire(this._getCurrentPageCode(), this.timer_minutes_total);
        }
      }
    },
  }
}
</script>
