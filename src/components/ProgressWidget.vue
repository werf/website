<template>
  <div class="progresswidget">
    <div v-if="skippedForm">
      <p>Для большей эффективности желательно <a :href="baseUrl+'/'+guideCode+'/000_task.html'">определиться с целями</a></p>
    </div>
    <div v-if="! skippedForm">
      <div class="progresswidget__target">Цель:
        <strong v-if="wish===BECOME_SPECIALIST">стать более востребованным специалистом</strong>
        <strong v-if="wish===BRING_KUBERNETES">принести Kubernetes в компанию</strong>
        <strong v-if="wish===SOLVE_TASKS">эффективнее решать задачи</strong>
      </div>
      <div class="progresswidget__progress">
        Выполнено шагов: <span>{{completedStepsCalculated}}</span> из {{learningProgress.length}}. Решение текущего шага: <div id="minutes_total" :class="'progresswidget__timer '+((timer_minutes_total>=timer_estimated_minutes)?'progresswidget__timer_overtime':'')">{{ timer_minutes_total }}</div> из
        {{ timer_estimated_minutes }} минут.
        &nbsp;&nbsp;
        <div :class="'button__blue button__blue_small button__blue_inline '+((timer_minutes_total>=timer_estimated_minutes)?'button__blue_red':'')">
          <a href="https://t.me/werf_ru">нужна помощь</a>
        </div>
        &nbsp;&nbsp;
        <div class="button__blue button__blue_small button__blue_inline" v-if="!expandedList">
          <a href="#" @click="expandList">свериться с планом</a>
        </div>
      </div>
      <div v-if="expandedList" class="learningprogress">
        <div class="learningprogress__checklist">
          <div v-for="el in learningProgress" :class="'learningprogress__row '+(el['completed']?'learningprogress__row_completed':'')">
            <span class="learningprogress__time">{{el['spent_minutes']}} из {{el['estimated_minutes']}}</span>
            <span class="learningprogress__label">{{el['label']}}</span>
          </div>
        </div>
        <div class="button__blue button__blue_small button__blue_inline">
          <a href="#" @click="hideList">свернуть план</a>
        </div>
        <div class="button__blue button__blue_small button__blue_inline">
          <a :href="baseUrl+'/'+guideCode+'/000_task.html'">изменить план</a>
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
        if (this._getCurrentPageTrackerMinutes() + 20 < this.timer_minutes_total) {
          this._setCurrentPageTrackerMinutes(this.timer_minutes_total);
          this._progresstracker_timer_fire(this._getCurrentPageCode(), this.timer_minutes_total);
        }
      }
    },
  }
}
</script>

<style lang="scss" scoped>
$gray: darkslategrey;
h1 {
  color: $gray;
  font-size: 2rem;
}
</style>
