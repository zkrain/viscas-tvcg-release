<template>
  <div class="infer-container" @mouseleave="showSlider = false">
    <div class="infer-btn" @click="infer" @mouseenter="showSlider = true">Infer</div>
    <transition name="slide-fade">
      <div v-if="showSlider" class="box">
        <div class="slider-box">
          <div class="params-label">K:</div>
          <el-slider class="slider-cls" v-model="k" :max="10" :min="0"></el-slider>
        </div>
        <div class="slider-box">
          <label class="params-label">TW:</label>
          <el-slider class="slider-cls" v-model="tw" :max="15" :min="1"></el-slider>
        </div>
      </div>
    </transition>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { Slider } from 'element-ui'
import '@/css/transition.css'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'

import EventBus from '@/eventBus'
const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    FontAwesomeIcon,
    ElSlider: Slider
  }
})
export default class Panel extends Vue {
  @projStore.Action private inferCascadingPattern!: ({ locationIDs, k, tw }:
                                                       { locationIDs: number[], k: number, tw: number }) => void
  @projStore.Getter private edgeLengthDict!: {[key: string]: number}
  private inferLocationIds: number[] = []
  private k: number = 0
  private tw: number = 6
  private showSlider: boolean = false

  private mounted (): void {
    EventBus.$on('setInferLocationIds', this.setInferLocationIds)
  }

  private setInferLocationIds (ids: number[]): void {
    this.inferLocationIds = ids
  }

  private infer (): void {
    if (this.inferLocationIds.length === 0) {
      alert('please select some points')
    } else {
      this.inferCascadingPattern({
        locationIDs: this.inferLocationIds,
        k: this.k === 0 ? 4 : this.k,
        tw: this.tw
      })
    }
  }
}
</script>

<style scoped lang="scss">
.infer-container {

  .infer-btn {
    background-color: white;
    border: 1px solid #eee;
    width: 100px;
    height: 26px;
    display: table-cell;
    text-align: center;
    vertical-align: middle;
    cursor: pointer;
    box-sizing: border-box;
    position: relative;
    border-radius: 4px;
    z-index: 999;
    transition: all 300ms;
    &:hover {
      color: #1cc5f0;
    }
  }

  .box {
    position: relative;
    margin-top: -16px;
    padding: 10px;
    width: fit-content;
    background-color: white;
    border: 1px solid #eee;
    margin-left: -29%;
    border-radius: 4px;

    .slider-box {
      width: 200px;
      height: 30px;
      display: flex;
      align-items: center;
      box-sizing: border-box;
      margin-bottom: 4px;

      .params-label {
        width: 20%;
        float: left;
      }

      .slider-cls {
        width: 70%;
        height: 100%;
        float: left;
      }
    }
  }
}
</style>
