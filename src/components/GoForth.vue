<template>
  <div v-if="!allCompleteCalculated">
    <div v-if="skippedForm">
      <a :href="url" class="nav-btn">Далее: {{label}}</a>
    </div>
    <div v-if="!skippedForm && nextTaskCalculated">
      <a :href="baseUrl+'/'+guideCode+'/'+nextTaskCalculated.url" class="nav-btn" @click="buttonClicked">Двигаемся дальше: {{nextTaskCalculated.menu_label}}</a>
    </div>
    <div v-if="!skippedForm && ! nextTaskCalculated">
      <a :href="baseUrl+'/'+guideCode+'/999_complete.html'" class="nav-btn" @click="finalButtonClicked">Отметить последнюю главу решённой</a>
    </div>
  </div>
</template>

<script>
import sharedMixin from "../shared"
import progressTrackerMixin from "../progresstracker"
export default {
  name: 'GoForth',

  mixins: [sharedMixin, progressTrackerMixin],

  props: {
    url: String,
    label: String,
  },

  data() {
    return {
    };
  },

  methods: {
    buttonClicked: function(event) {
      this.markCurrentPageComplete(event);
      this._progresstracker_complete_step(this._getCurrentPageCode());
    },

    finalButtonClicked: function (event) {
      this.markCurrentPageComplete(event);
      this._progresstracker_complete_step(this._getCurrentPageCode());
      this._progresstracker_complete_all();
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
