<template>
  <div class="main">
    <Menu />
    <div class="content">
      <div class="left">
        <div class="left-content">
          <div class="top">
            <map-view />
          </div>
          <div class="bottom">
            <explanation-view />
          </div>
        </div>
      </div>
      <div class="right">
        <div class="tree-statistics">
          <tree-statistics-view />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import ExplanationView from './ExplanationView/index.vue'
import MapView from './MapView/index.vue'
import TreeStatisticsView from './TreeStatisticsView/index.vue'
import Menu from './Menu/index.vue'

import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    MapView,
    ExplanationView,
    Menu,
    TreeStatisticsView
  }
})
export default class Main extends Vue {
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.main {
  position: relative;
  width: 100%;
  height: 100%;
  flex-direction: column;
  display: flex;

  .content {
    position: relative;
    width: 100%;
    height: calc(100% - 40px);
    background:rgba(247,247,247,1);

    .left {
      position: relative;
      height: 100%;
      // width: calc(100% - 536px);
      width: calc(100% - 602px);
      border-right: 1px solid #eee;
      float: left;

      .left-content {
        width: calc(100% - 12px);
        height: 100%;
        position: relative;

        .top {
          position: relative;
          width: calc(100% - 1px);
          height: calc(100% - 601px);
          border-bottom: 1px solid #eee;
          border-right: 1px solid #eee;
          box-shadow: 1px 1px 5px 1px rgba(201,201,201,0.4);
        }

        .bottom {
          position: relative;
          width: 100%;
          height: calc(600px - 12px);
          margin-top: 12px;
          box-shadow: 1px -1px 5px 1px rgba(201,201,201,0.4);
        }
      }
    }

    .right {
      position: relative;
      height: 100%;
      // width: 535px;
      width: 600px;
      float: left;
      background-color: #fff;
      box-shadow: -1px 0px 5px 1px rgba(201,201,201,0.4);

      .tree-statistics {
        position: relative;
        flex: 0 0 600px;
        height: 100%;
        overflow: hidden;

        .edge-handles {
          position: absolute;
          left: 1px;
          top: 0px;
          width: calc(100% - 2px);
          height: 100%;
          display: flex;
          pointer-events: none;

          .edge-handle {
            flex: 1;
            border: 2px solid rgba(0,0,0,0.0);
            transition: border 300ms;
            border-radius: 2px;
          }

          .edge-handle:hover {
            border: 2px solid rgb(28, 197, 240);
          }

          .selected {
            border: 2px solid rgb(28, 197, 240) !important;
          }
        }
      }
    }
  }
}
</style>
