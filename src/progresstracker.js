export default {
  methods: {
    _progresstracker_complete_form(wish, learningProgress) {
      ym(67759933,'reachGoal','complete_form');
      ym(67759933,'reachGoal','audience_' + wish);
      for (let i in learningProgress) {
        ym(67759933,'reachGoal',learningProgress[i]['code']+'_open');
      };
    },

    _progresstracker_complete_step(code){
      ym(67759933,'reachGoal',code+'_resolved');
    },

    _progresstracker_complete_all(){
      ym(67759933,'reachGoal','complete_guide_open');
    },

    _progresstracker_timer_fire(code, time){
      ym(67759933,'reachGoal',code + '_timer_' + time);
    },
  }
}
