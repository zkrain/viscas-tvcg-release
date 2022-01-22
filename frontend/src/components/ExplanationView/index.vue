<template>
  <div class="explanation-view">
    <div class="tree-pattern">
      <div class="trees-cont">
        <dynamic-trees-view/>
      </div>
      <div class="cascades">
<!--        <div class="locationID">-->
<!--          <div class="id-box">-->
<!--            <div v-for="(val, key) in showLocationIds" :key="key" class="id-cell">{{val}}</div>-->
<!--          </div>-->
<!--        </div>-->
        <div class="time-windows" ref="timeWindows" @scroll="onScroll($event)">
          <transition-group name="list-complete">
            <time-line-box v-for="r in showTimeWindows" :key="r.tree.id"
                           :time-window-info="r" :locations-id="showLocationIds"/>
          </transition-group>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import DynamicTreesView from './DynamicTreesView.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import { uniq } from 'lodash'
import TimeLineBox from '@/components/ExplanationView/TimeLineBox.vue'

const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    TimeLineBox,
    FontAwesomeIcon,
    DynamicTreesView
  }
})
export default class ExplanationView extends Vue {
  @projStore.State private timeWindows!: ITimeWindowInfo[]
  @projStore.State private cascadingLinks!: number[][]
  private showTimeWindows: ITimeWindowInfo[] = []

  private onScroll (e: Event): void {
    const target: HTMLElement = e.currentTarget as HTMLElement
    const dy: number = target.scrollWidth - target.clientWidth - target.scrollLeft
    if (dy < target.clientWidth / 2 && this.showTimeWindows.length < this.timeWindows.length) {
      console.log('lazy load.')
      const nowPos: number = this.showTimeWindows.length
      for (let i = nowPos; i < this.timeWindows.length && i < nowPos + 10; i++) {
        this.showTimeWindows.push(this.timeWindows[i])
      }
    }
  }

  @Watch('timeWindows')
  private setShowTimeWindows (): void {
    this.showTimeWindows = []
    if (typeof this.timeWindows === 'undefined') {
      return
    }
    if (this.timeWindows.length <= 20 && typeof this.timeWindows !== 'undefined') {
      this.timeWindows.forEach((t) => {
        this.showTimeWindows.push(t)
      })
    } else {
      for (let i = 0; i < 20; i++) {
        this.showTimeWindows.push(this.timeWindows[i])
      }
    }
  }

  get showLocationIds (): number[] {
    const locationIds: number[] = []
    this.cascadingLinks.forEach((link: number[]) => {
      locationIds.push(link[0], link[1])
    })
    return uniq(locationIds)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.explanation-view {
  position: relative;
  border: 1px solid #eee;
  border-radius: 4px;
  height: 100%;
  width: 100%;
  user-select: none;
  background-color: white;

  .tree-pattern {
    position: relative;
    height: 100%;
    width: 100%;

    .trees-cont {
      position: relative;
      height: 400px;
      width: calc(100% - 1px);
    }

    .cascades {
      position: relative;
      height: calc(100% - 414px);
      width: calc(100% - 1px);
      box-sizing: border-box;
      padding-left: 16px;
      padding-right: 16px;
      background-color: white;
      margin-top: 12px;
      border-radius: 4px;

      .time-windows {
        position: relative;
        float: left;
        height: 100%;
        width: 100%;
        box-sizing: border-box;
        overflow-x:scroll;
        overflow-y:hidden;
        white-space: nowrap;
        text-align: left;
      }
    }
  }
}
.timeline-slider {
  height: 300px;
  width: 100%;
  z-index: 1;
}
::-webkit-scrollbar {
  display: none;
}

.list-complete-enter, .list-complete-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
.list-complete-leave-active {
  transition: left 0ms;
}
</style>
