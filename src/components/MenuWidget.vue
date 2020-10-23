<template>
</template>

<script>
import sharedMixin from "../shared"
import progressTrackerMixin from "../progresstracker"
export default {
  name: 'MenuWidget',

  mixins: [sharedMixin, progressTrackerMixin],

  data() {
    return {
      expandedList: false
    };
  },

  created() {
    this._updateMenuItems();
  },

  methods: {
    _updateMenuItems(){
      if (this.skippedForm)
        return;

      let menu_items_all = document.querySelectorAll('[data-subhref]');
      menu_items_all.forEach(function(element) {
        if (! element.parentElement.classList.contains('sidebar__item_parent'))
          element.setAttribute('class','unselected');
      })

      if (this.$data.learningProgress)
        for (let i in this.learningProgress) {
          let step = this.learningProgress[i];
          let menu_items_highlighted = document.querySelectorAll('[data-subhref="'+step.url+'"]');
          menu_items_highlighted.forEach(function(element){
            if (step.completed) {
              element.setAttribute('class','completed');
            }else{
              element.setAttribute('class','todo');
            }
          })
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
