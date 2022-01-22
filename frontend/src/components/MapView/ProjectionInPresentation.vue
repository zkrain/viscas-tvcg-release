<template>
  <div class="projection">
    <div class="projection-canvas" ref="svgContRef">
      <svg :width="svgWidth" :height="svgHeight">
        <circle v-for="p in embedPoints" :key="p.rid" :id="p.rid"
                :cx="p.left * svgWidth" :cy="p.top * svgHeight"
                :r="filteredIDDict[p.rid] ? 3 : 1"
                :fill="colorMapForEachLocation[p.rid] ? colorMapForEachLocation[p.rid]: '#ccc'"/>
        <rect x="0" y="0" :width="svgWidth" :height="svgHeight" v-if="showMask" fill="white" opacity="0.8"/>
      </svg>
      <svg :width="svgWidth" :height="svgHeight" class="cluster-circle-cls">
        <defs>
          <filter id="blurFilter" x="-5" y="-5" width="40" height="40">
            <feGaussianBlur in="SourceAlpha" result="blurOut" stdDeviation="4"/>
            <feColorMatrix in="blurOut" type="matrix" result="colorOut"
                           values="1 1 1 0.8 0
                                    1 1 1 0.8 0
                                    1 1 1 0.8 0
                                    1 1 1 0.8 0"/>
            <feMerge>
              <feMergeNode in="colorOut" />
              <feMergeNode in="SourceGraphic" />
            </feMerge>
          </filter>
        </defs>
        <g v-for="(p, key) in circlesInfo" :key="key"
           :transform="`translate(${p.center[0]},${p.center[1]}) scale(${p.clusterId === selectedCluster? 1.2: 1})`">
          <circle :r="p.radius" :filter="p.clusterId === selectedCluster? 'url(#blurFilter)' : 'none'"
                  @click="onClusterClick(p.rids, p.clusterId)"
                  class="cluster-circle" />
          <circle :r="p.radius - 6" :opacity="0.5"
                  @click="onClusterClick(p.rids, p.clusterId)"
                  class="cluster-circle" />
          <circle :cx="p.drawCenterPosition[0]" :cy="p.drawCenterPosition[1]" :r="smallCircleRadius"
                  class="small-circle"/>
          <path :d="p.boxPlotPath" class="cluster-arc box"/>
        </g>
        <!--重复绘制来解决一下压盖问题，没办法了，只能这样子-->
        <g v-for="(p, key) in circlesInfo" :key="'h' + key">
          <g v-if="p.clusterId === selectedCluster"
             :transform="`translate(${p.center[0]},${p.center[1]}) scale(${p.clusterId === selectedCluster? 1.2: 1})`">
            <circle :r="p.radius" :filter="p.clusterId === selectedCluster? 'url(#blurFilter)' : 'none'"
                    @click="onClusterClick(p.rids, p.clusterId)"
                    class="cluster-circle" />
            <circle :r="p.radius - 6" :opacity="0.5"
                    @click="onClusterClick(p.rids, p.clusterId)"
                    class="cluster-circle" />
            <circle :cx="p.drawCenterPosition[0]" :cy="p.drawCenterPosition[1]" :r="smallCircleRadius"
                    class="small-circle"/>
            <path :d="p.boxPlotPath" class="cluster-arc box"/>
          </g>
        </g>
      </svg>
    </div>
    <div class="projection-panel">
      <div class="parameters-cont">
        <div class="param-name">EpsT:</div>
        <div class="input-cont">
          <el-slider v-model="EpsT" :show-tooltip="false" :max="0.1" :min="0" :step="0.001"/>
        </div>
        <div class="value-cont">{{EpsT}}</div>
      </div>
      <div class="parameters-cont">
        <div class="param-name">EpsS:</div>
        <div class="input-cont">
          <el-slider v-model="EpsS" :show-tooltip="false"
                     :max="currentDataSetConfig.rangeEpsS[1]"
                     :min="currentDataSetConfig.rangeEpsS[0]" :step="0.1"/>
        </div>
        <div class="value-cont">{{EpsS}}</div>
      </div>
      <div class="parameters-cont">
        <div class="param-name">#Pts:</div>
        <div class="input-cont">
          <el-slider v-model="MinPts" :show-tooltip="false" :max="25" :min="0"/>
        </div>
        <div class="value-cont">{{MinPts}}</div>
      </div>
      <div class="parameters-cont">
        <div class="btn-cont">
          <div class="clustering-btn" @click="triggerClustering">cluster</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import EventBus from '../../eventBus'

import { distance } from '@/utils/geoCompute'
import { arcAngleToPath, arcAngleRangePathD, getBoxPlotInfo } from '@/utils/draw'
import { Slider } from 'element-ui'

import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    ElSlider: Slider
  }
})
export default class ProjectionInPresentation extends Vue {
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.Action private getRoadsByID!: (roadsID: number[]) => void
  @projStore.Action private setProjectionCanvasSize!: (payload: {width: number, height: number}) => void
  @projStore.Action private getRoadsSample!: (roadsIDSample: number[]) => void
  @projStore.Action private getCorrelatedLocations!: (params: { clusterId: number, locationIDs: number[]}) => void
  @projStore.Action private stclustering!: (payload: {
    EpsT: number,
    EpsS: number,
    MinPts: number,
    deltaE: number
  }) => { clusters: number[][], noise: number[] }
  @projStore.State private embedPoints!: ProjPoint[]
  @projStore.State private clustersWithContext!: ICluster[]
  @projStore.State private colorMapForEachLocation!: {[rid: number]: string}
  @projStore.Getter private filteredIDDict!: {[key: number]: number}

  private EpsT = 0.085 // distance in projection
  private EpsS = 2.5 // km
  private MinPts = 10
  private deltaE = 0.003 // distance in projection

  private selectedCluster = -1
  private circlesInfo: ICircleInfo[] = []
  private circleRadius: number = 24
  private smallCircleRadius: number = 3
  private svgWidth: number = 416
  private svgHeight: number = 358
  private showMask: boolean = false

  private mounted (): void {
    this.setProjectionCanvasSize({ width: this.svgWidth, height: this.svgHeight })
  }

  private async triggerClustering () { // ST-DBSCAN
    // this.removeAllMaps()
    // // ******** computation in backend ********
    this.showMask = true
    this.stclustering({
      EpsT: this.EpsT,
      EpsS: this.EpsS,
      MinPts: this.MinPts,
      deltaE: this.deltaE
    })
  }

  private onClusterClick (locationsId: number[], clusterId: number): void {
    this.selectedCluster = clusterId
    EventBus.$emit('onClickCluster', locationsId, clusterId)
  }

  @Watch('currentDataSetConfig')
  private onChangeDataSet (): void {
    this.circlesInfo = []
    this.showMask = false
    this.selectedCluster = -1
    this.EpsT = this.currentDataSetConfig.initParams[0]
    this.EpsS = this.currentDataSetConfig.initParams[1]
    this.MinPts = this.currentDataSetConfig.initParams[2]
  }

  @Watch('clustersWithContext')
  private onclustersWithContextChanged (clustersWithContext: ICluster[]) {
    this.initCircleInfo()
  }

  private initCircleInfo (): void {
    this.circlesInfo = []
    const tmpCircleInfo: ICircleInfo[] = []
    let maxAverageEventCount: number = 0
    let minAverageEventCount: number = 99999
    for (const tmpCluster of this.clustersWithContext) {
      maxAverageEventCount = Math.max(maxAverageEventCount, tmpCluster.clusterAverageEventCount)
      minAverageEventCount = Math.min(minAverageEventCount, tmpCluster.clusterAverageEventCount)
    }
    let averageEventCountK: number = 1
    if (maxAverageEventCount > minAverageEventCount) {
      averageEventCountK = 1 / (maxAverageEventCount - minAverageEventCount)
    }
    let eventCountsK: number = 1
    let minEventCounts: number = 99999
    let maxEventCounts: number = 0
    this.clustersWithContext.forEach((tmpCluster: ICluster) => {
      tmpCluster.clusterEventCount.forEach((v: number) => {
        minEventCounts = Math.min(v, minEventCounts)
        maxEventCounts = Math.max(v, maxEventCounts)
      })
    })
    eventCountsK = 1 / (maxEventCounts - minEventCounts)
    this.clustersWithContext.forEach((tmpCluster: ICluster, key: number) => {
      const tmpHullPoints = tmpCluster.hullPoints
      let sumX: number = 0
      let sumY: number = 0
      tmpHullPoints.forEach((hullPoint: [number, number]) => {
        sumX += hullPoint[0]
        sumY += hullPoint[1]
      })
      let rx: number = this.circleRadius * (tmpCluster.spatialRelativeCenter.rx >= 0
        ? Math.pow(tmpCluster.spatialRelativeCenter.rx, 0.7) : -Math.pow(-tmpCluster.spatialRelativeCenter.rx, 0.7))
      let ry: number = -this.circleRadius * (tmpCluster.spatialRelativeCenter.ry >= 0
        ? Math.pow(tmpCluster.spatialRelativeCenter.ry, 0.7) : -Math.pow(-tmpCluster.spatialRelativeCenter.ry, 0.7))
      const arcAngleRate: number = averageEventCountK * (tmpCluster.clusterAverageEventCount - minAverageEventCount) // 标准化
      // 小圆边界超出中心的处理, 1 表示容差， 可调
      if (Math.sqrt(rx * rx + ry * ry) + this.smallCircleRadius > this.circleRadius + 1) {
        let dr: number = this.circleRadius + 1 - this.smallCircleRadius - Math.sqrt(rx * rx + ry * ry)
        rx = rx - dr * rx / Math.sqrt(rx * rx + ry * ry)
        ry = ry - dr * ry / Math.sqrt(rx * rx + ry * ry)
      }
      let center: [number, number] = [sumX / tmpHullPoints.length, sumY / tmpHullPoints.length]
      const boxInfo: number[] = getBoxPlotInfo(tmpCluster.clusterEventCount)
      let minAng: number = 2 * Math.PI * eventCountsK * (boxInfo[1] - minEventCounts)
      let maxAng: number = 2 * Math.PI * eventCountsK * (boxInfo[3] - minEventCounts)
      tmpCircleInfo.push({
        clusterId: key,
        rids: tmpCluster.rIDs,
        color: tmpCluster.color,
        arcPath: arcAngleToPath(arcAngleRate * 2 * Math.PI, [0, 0], this.circleRadius - 3),
        boxPlotPath: arcAngleRangePathD([minAng, maxAng], [0, 0], this.circleRadius - 3),
        center: center,
        radius: this.circleRadius,
        hullPoints: tmpCluster.hullPoints,
        drawCenterPosition: [rx, ry]
      })
    })
    // 偏移以消除重叠部分
    tmpCircleInfo.forEach(tCircleInfo => {
      const intersectCircle = []
      const moveDirection: [number, number] = [0, 0]
      for (const circleInfo of this.circlesInfo) {
        if (tCircleInfo.radius + circleInfo.radius > distance(tCircleInfo.center, circleInfo.center)) {
          intersectCircle.push(circleInfo)
          moveDirection[0] += (tCircleInfo.center[0] - circleInfo.center[0])
          moveDirection[1] += (tCircleInfo.center[1] - circleInfo.center[1])
        }
      }
      if (intersectCircle.length <= 0) {
        this.circlesInfo.push(tCircleInfo)
      } else { // 如果效率比较低的话，这里或许可以优化
        const moveDirectionLen = Math.sqrt(moveDirection[0] * moveDirection[0] + moveDirection[1] * moveDirection[1])
        moveDirection[0] /= moveDirectionLen
        moveDirection[1] /= moveDirectionLen
        let stopMove: boolean = false
        const maxStep: number = 1000 // 最大移动步数
        let step: number = 0
        while (!stopMove && step <= maxStep) {
          tCircleInfo.center[0] += moveDirection[0]
          tCircleInfo.center[1] += moveDirection[1]
          let flagStopMove: boolean = true
          for (const circleInfo of this.circlesInfo) {
            if (tCircleInfo.radius + circleInfo.radius > distance(tCircleInfo.center, circleInfo.center)) {
              flagStopMove = false
              break
            }
          }
          stopMove = flagStopMove
          step += 1
        }
        this.circlesInfo.push(tCircleInfo)
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.projection {
  position: relative;
  width: 100%;
  height: 100%;
  border-right: 1px solid #eee;

  .cluster-circle-cls {
    top: calc(-100% - 8px);
    z-index: 3;
  }

  .projection-panel {
    position: relative;
    width: 100%;
    height: 76px;
    border-top: 1px solid #eee;

    .parameters-cont {
      position: relative;
      float: left;
      height: 38px;
      width: calc(50% - 20px);
      margin-left: 20px;
      display: flex;

      .btn-cont {
        position: relative;
        height: 100%;
        width: 100%;

        .clustering-btn {
          position: relative;
          height: 26px;
          border: 1px solid #eee;
          width: calc(100% - 20px);
          margin: 6px 0;
          border-radius: 2px;
          cursor: pointer;
          line-height: 26px;

          &:hover {
            color: #1cc5f0;
          }
        }
      }

      .param-name {
        margin-right: 10px;
        line-height: 38px;
        flex: 0 0 38px;
        height: 100%;
        text-align: start;
        font-size: 14px;
      }

      .input-cont {
        margin-right: 10px;
        height: 100%;
        line-height: 38px;
        flex: 1;
      }

      .value-cont {
        flex: 0 0 45px;
        text-align: start;
        line-height: 38px;
      }
    }
  }

  .projection-canvas {
    position: relative;
    width: 100%;
    height: calc(100% - 77px);
    border: 1px solid #eee;
    border-bottom: none;

    .div-point {
      position: absolute;
      border-radius: 50%;
      height: 3px;
      width: 3px;
      background-color: #ccc;
    }

    svg {
      position: relative;
    }

    .not_possible {
      fill: rgb(200,200,200);
    }

    .possible {
      fill: #EC888C;
    }

    .selected {
      fill: steelblue;
      stroke-width: 1;
      stroke: #fff;
    }

    .cluster-circle {
      fill: white;
      cursor: pointer;
      stroke-width: 1px;
      stroke: #aaa;
    }

    .small-circle {
      fill: none;
      stroke-width: 1px;
      stroke: #aaa;
    }

    .cluster-arc {
      stroke-width: 6px;
      stroke: orange;
      opacity: 0.5;
      fill: none;
    }

    .box {
      stroke: #777;
      opacity: 1;
    }
  }
}
</style>
