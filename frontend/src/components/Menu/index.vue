<template>
  <div class="menu-bar">
    <div class="btn-dataset-cont">
      <div class="select-btn">
        <div class="selected-data-set-name">{{currentDataSetConfig.name}}</div>
        <div class="icon" @click="triggerDatasetSelection">
          <font-awesome-icon :icon="selectingDataset? faAngleUpIcon : faAngleDownIcon"/>
        </div>
      </div>
    </div>
    <transition name="slide-fade">
      <div v-if="selectingDataset" class="datasets">
          <div class="dataset" v-for="dataset in datasets" :key="dataset.name"
               @click="selectDataset(dataset)">{{dataset.name}}</div>
      </div>
    </transition>

    <div class="title">
      <div class="icon">
        <font-awesome-icon :icon="faProjectDiagramIcon" />
      </div>
      <div class="name">
        VisCas
      </div>
    </div>
    <div class="Tutorial fake">
      <div class="container">
        <div class="icon">
          <font-awesome-icon :icon="faStickyNoteIcon" />
        </div>
        <div class="name">
          Tutorial
        </div>
      </div>
    </div>
    <div class="About fake">
      <div class="container">
        <div class="icon">
          <font-awesome-icon :icon="faInfoIcon" />
        </div>
        <div class="name">
          About
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import 'animate.css'
import '@/css/transition.css'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faProjectDiagram, faAngleDown, faAngleUp, faStickyNote, faInfo } from '@fortawesome/free-solid-svg-icons'

import EventBus from '../../eventBus'

const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    FontAwesomeIcon
  }
})
export default class Menu extends Vue {
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.Getter private locationRange!: ILocationRange
  @projStore.Action private setDataSetConfig!: (dataSetConfig: IDataSetConfig) => void
  @projStore.Action private fetchDataset!: (filenameAbbrev: string) => void
  @projStore.Action private setPieChartInnerRate!: (rate: number) => void
  private selectingDataset: boolean = false

  private datasets: IDataSetConfig[] = [{
    name: 'Congestion',
    abbrev: 'congestion2',
    mapCenter: [30.3, 120.15],
    zoom: 14,
    rangeEpsS: [0, 6],
    initParams: [0.085, 2.5, 10],
    pieChartInnerRate: 8,
    donutChartSmallRadius: 4,
    donutSectorNum: 24,
    glyphMinZoom: 14,
    secondGlyphZoom: 15,
    startTime: new Date(2018, 2, 1, 0, 0, 0),
    deltaTime: 10 * 60,
    xTimeBin: 61 + 1,
    yTimeBin: 24 * 6,
    yTimeBinRate: 6,
    gridShape: [61 + 1, 24],
    totalInfluenceDuration: 16,
    pLowerBound: 0.5,
    timelineLengthRate: 300 / 18
  }, {
    name: 'UrbanFlow',
    abbrev: 'flow',
    mapCenter: [39.907, 116.383331],
    zoom: 12,
    rangeEpsS: [0, 6],
    initParams: [0.066, 2.5, 10],
    startTime: new Date(2015, 10, 1),
    deltaTime: 30 * 60,
    pieChartInnerRate: 3,
    donutChartSmallRadius: 4,
    donutSectorNum: 24,
    glyphMinZoom: 14,
    secondGlyphZoom: 15,
    xTimeBin: 161 + 1,
    yTimeBin: 24 * 2,
    yTimeBinRate: 2,
    gridShape: [161 + 1, 24],
    totalInfluenceDuration: 20,
    pLowerBound: 0.2,
    timelineLengthRate: 20
  }, {
    name: 'AirPollution',
    abbrev: 'air',
    mapCenter: [38, 105],
    zoom: 4,
    rangeEpsS: [100, 2000],
    initParams: [0.026, 670, 10],
    startTime: new Date(2018, 0, 1),
    deltaTime: 60 * 60,
    pieChartInnerRate: 1,
    donutChartSmallRadius: 6,
    donutSectorNum: 12,
    glyphMinZoom: 7,
    secondGlyphZoom: 8,
    xTimeBin: 50 + 1,
    yTimeBin: 24 * 7,
    yTimeBinRate: 7,
    gridShape: [50 + 1, 24],
    totalInfluenceDuration: 60,
    pLowerBound: 0.2,
    timelineLengthRate: 200 / 18
  }]

  private async created () {
    await this.fetchDataset(this.currentDataSetConfig.abbrev)
    this.setPieChartInnerRate(this.currentDataSetConfig.pieChartInnerRate)
  }

  private mounted (): void {
    EventBus.$emit('updateMapCenterAndZoom', this.currentDataSetConfig.mapCenter, this.currentDataSetConfig.zoom)
  }

  get faAngleDownIcon () {
    return faAngleDown
  }

  get faAngleUpIcon () {
    return faAngleUp
  }

  get faProjectDiagramIcon () {
    return faProjectDiagram
  }

  get faStickyNoteIcon () {
    return faStickyNote
  }

  get faInfoIcon () {
    return faInfo
  }

  private triggerDatasetSelection () {
    this.selectingDataset = !this.selectingDataset
  }

  private selectDataset (dataset: IDataSetConfig) {
    this.selectingDataset = !this.selectingDataset
    this.setDataSetConfig(dataset)
    this.fetchDataset(dataset.abbrev)
    this.setPieChartInnerRate(dataset.pieChartInnerRate)
    EventBus.$emit('updateMapCenterAndZoom', dataset.mapCenter, dataset.zoom)
  }
}
</script>

<style lang="scss" scoped>
.menu-bar {
  position: relative;
  flex: 0 0 40px;
  height: 40px;
  width: 100%;
  background:rgba(5,30,62,1);
  box-shadow:0px 2px 4px 0px rgba(32,84,159,0.4);
  z-index: 2;

  .btn-dataset-cont {
    position: absolute;
    top: 7px;
    left: 200px;
    height: 24px;
    width: 150px;

    .select-btn {
      position: relative;
      height: 100%;
      width: 100%;
      display: flex;
      background-color: #fff;
      border-radius: 3px;
      border: 1px solid #ddd;
      line-height: 24px;

      .selected-data-set-name {
        color: #2c3e50;
        flex: 1;
      }

      .icon {
        flex: 0 0 30px;
        border-left: 1px solid #eee;
        cursor: pointer;
        transition: all 300ms;
      }
    }
  }

  .datasets {
    position: absolute;
    top: 36px;
    left: 200px;
    width: 120px;
    height: 60px;
    background-color: #fff;
    border-radius: 2px;
    border: 1px solid #eee;
    z-index: 100;

    .dataset {
      position: relative;
      width: calc(100% - 10px);
      height: 20px;
      background-color: #fff;
      line-height: 20px;
      text-align: start;
      font-size: 16px;
      transition: all 300ms;
      padding-left: 10px;
      cursor: pointer;
    }

    .dataset:hover {
      background-color: #3399FF;
      color: #fff;
    }
  }

  .title {
    position: absolute;
    left: 20px;
    top: 0px;
    height: 100%;
    width: 300px;
    line-height: 40px;

    .icon {
      position: relative;
      float: left;
      width: 50px;
      color: white;
      font-size: 22px;
    }

    .name {
      position: relative;
      float: left;
      width: calc(100% - 50px);
      color: #2c3e50;
      text-align: start;
      font-size:20px;
      font-family:Helvetica-Bold,Helvetica,serif;
      font-weight:bold;
      color:rgba(255,255,255,1);
    }
  }

  .fake {
    position: absolute;
    width: 100px;
    height: 100%;

    .container {
      position: relative;
      height: 22px;
      top: 9px;
      background-color: #fff;
      border-radius: 2px;
      line-height: 22px;
      color: #2c3e50;
      font-size: 18px;
      border: 1px solid #eee;

      .icon {
        position: relative;
        float: left;
        width: 30px;
        top: 2px;
      }

      .name {
        position: relative;
        float: left;
        width: calc(100% - 30px);
        text-align: start;
      }
    }
  }

  .Tutorial {
    left: 85%;
  }

  .About {
    left: calc(85% + 120px);
  }

  .projection {
    left: 63%;
    user-select: none;
  }
}
</style>
