<template>
  <div class="dynamic-trees-view">
    <div class="title-div">
      <div class="name">Cascading View</div>
    </div>
    <div class="dynamic-trees-view-cont">
      <div class="top-div">
        <svg class="head-svg">
          <rect v-for="v in xTicks" :key="`${v}-head`" :x="getPositionAlongX(v * yTimeBin)" y="0"
                :width="getPositionAlongX(yTimeBin)" class="cell" :opacity="columnsSum[v]"/>
        </svg>
      </div>
      <div class="mid-div">
        <div class="trees-timeline">
          <div class="timeline" ref="timelineRef">
            <svg class="trees-svg" ref="treesSvgRef">
              <g ref="gridG">
                <line v-for="v in xTicks" :key="`${v}-v`" class="line-vertical"
                      :x1="getPositionAlongX(v * yTimeBin)"
                      :x2="getPositionAlongX(v * yTimeBin)"
                      y1="0" :y2="timelineContainerHeight" style="stroke:#ddd;stroke-width:0.5" />
                <line v-for="v in yTicks" :key="`${v}-h`" class="line-horizon"
                      :y1="getPositionAlongY(v)"
                      :y2="getPositionAlongY(v)"
                      x1="0" :x2="timelineContainerWidth" style="stroke:#ddd;stroke-width:0.5" />
              </g>
            </svg>
            <div class="brusher"
                 @mousedown.left="onDragStarted($event)"
                 @mousemove.left="onDragUpdated($event)">
              <div v-if="selection.active || selection.animated"
                   :class="[{ animated: selection.animated } ,'selection']" :style="selectionStyle"></div>
            </div>
          </div>
        </div>
        <div class="right-div">
          <svg class="right-svg">
            <rect v-for="v in yTicks" :key="`${v}-right`" x="0" :y="getPositionAlongY(v)" class="cell"
                  :height="getPositionAlongY(currentDataSetConfig.yTimeBinRate)"
                  :opacity="rowsSum[v/currentDataSetConfig.yTimeBinRate]"/>
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import * as d3 from 'd3'
import Color from 'color'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'

const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    FontAwesomeIcon
  }
})
export default class DynamicTreesView extends Vue {
  @projStore.State private cascadingLinks!: number[][]
  @projStore.State private trees!: PropagationTree[]
  @projStore.State private orderedTrees!: PropagationTree[]
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.State private cascadingPointsScreenPosition!: {[locationID: number]: [number, number]}
  @projStore.Getter private edgeColorDict!: {[edgeKey: string]: [string, string]}
  @projStore.Getter private locationColorDict!: {[locationID: number]: string}
  @projStore.State private selectedTreesId!: number[]
  @projStore.State private selectedAggregatePropagationTreesId!: number[]
  @projStore.Getter private selectedTreesIdDictByATrees!: {[key: number]: boolean}
  @projStore.State private selectedEdge!: {[edgeKey: string]: boolean}
  @projStore.Getter private filteredTrees!: PropagationTree[]
  @projStore.Getter private filteredTreesIndices!: {[index: number]: boolean}
  @projStore.Action private dispatchTimeWindows!: (timeWindows: ITimeWindowInfo[]) => void

  private timelineContainerWidth = 1
  private timelineContainerHeight = 1
  private containerPageLeft = 0
  private containerPageTop = 0
  private patternG!: d3.Selection<SVGGElement, unknown, any, any>
  private rowsSum: number[] = []
  private columnsSum: number[] = []

  private selection = {
    animated: false,
    active: false,
    x1: 0,
    y1: 0,
    x2: 0,
    y2: 0
  }

  get selectionBound () {
    return {
      x: Math.min(this.selection.x1, this.selection.x2),
      y: Math.min(this.selection.y1, this.selection.y2),
      width: Math.abs(this.selection.x2 - this.selection.x1),
      height: Math.abs(this.selection.y2 - this.selection.y1)
    }
  }

  get selectionStyle () {
    return this.selection.animated ? this.animatedStyle : {
      left: this.selectionBound.x + 'px',
      top: this.selectionBound.y + 'px',
      width: this.selectionBound.width - 2 + 'px',
      height: this.selectionBound.height + 'px',
      border: `1px solid ${
        Color('#54a0ff').alpha(0.7).darken(0.1)}`,
      'background-color': Color('#54a0ff').alpha(0.3).toString(),
      'border-radius': 4 + 'px'
    }
  }

  @Watch('cascadingLinks')
  private resetSelection () {
    this.selection = {
      animated: false,
      active: false,
      x1: 0,
      y1: 0,
      x2: 0,
      y2: 0
    }
  }

  @Watch('selectedTreesId')
  private onSelectedTreesChanged (): void {
    this.patternG.selectAll('*').remove()
    this.drawRect(this.trees)
    this.resetTimeWindows(this.trees)
  }

  private initRowColumnSum (): void {
    this.columnsSum = new Array(this.xTicks.length + 1).fill(0)
    this.rowsSum = new Array(this.yTicks.length + 1).fill(0)
    this.trees.forEach((t: PropagationTree) => {
      const rowOk: boolean = (this.selectedTreesIdDictByATrees[t.id] && this.selectedAggregatePropagationTreesId.length > 0) ||
        (this.selectedAggregatePropagationTreesId.length === 0)
      const columnOk: boolean = t.id in this.filteredTreesIndices
      const startX: number = Math.floor(t.TreeStartTime / this.yTimeBin)
      const endX: number = Math.floor(t.TreeEndTime / this.yTimeBin)
      const startY: number = Math.floor((t.TreeStartTime % this.yTimeBin) / this.currentDataSetConfig.yTimeBinRate)
      const endY: number = Math.floor((t.TreeEndTime % this.yTimeBin) / this.currentDataSetConfig.yTimeBinRate)
      const maxY: number = this.yTimeBin / this.currentDataSetConfig.yTimeBinRate - 1
      if (rowOk && columnOk) {
        if (startX === endX) {
          this.columnsSum[startX] += (endY - startY + 1)
        } else {
          this.columnsSum[startX] += (maxY - startY + 1)
          this.columnsSum[endX] += (endY + 1)
        }
        if (startX === endX) {
          let t = startY
          while (t <= endY) {
            this.rowsSum[t] += 1
            t++
          }
        } else {
          let t = startY
          while (t < this.yTicks.length) {
            this.rowsSum[t] += 1
            t++
          }
          t = 0
          while (t <= endY) {
            this.rowsSum[t] += 1
            t++
          }
        }
      }
    })
    // normalized
    const maxRowSum: number = Math.max(...this.rowsSum)
    const minRowSum: number = Math.min(...this.rowsSum)
    if (maxRowSum > 0) {
      this.rowsSum = this.rowsSum.map((s) => {
        return (s - minRowSum) / (maxRowSum - minRowSum)
      })
    }
    const maxColumnSum: number = Math.max(...this.columnsSum)
    const minColumnSum: number = Math.min(...this.columnsSum)
    if (maxColumnSum > 0) {
      this.columnsSum = this.columnsSum.map((s) => {
        return (s - minColumnSum) / (maxColumnSum - minColumnSum)
      })
    }
  }

  private onDragStarted (e: MouseEvent) {
    this.selection.active = true
    this.selection.animated = false
    this.selection.x1 = this.selection.x2 = e.pageX
    this.selection.y1 = this.selection.y2 = e.pageY

    const stopFn = () => {
      document.removeEventListener('mouseup', stopFn)
      this.selection.active = false
      this.selection.animated = true
      this.retreiveBrushedTimeWindows()
    }
    document.addEventListener('mouseup', stopFn)
  }

  private onDragUpdated (e: MouseEvent) {
    if (this.selection.active) {
      this.selection.x2 = e.pageX
      this.selection.y2 = e.pageY
    }
  }

  get brushedTimeBin () {
    let x1 = this.selectionBound.x - this.containerPageLeft
    let y1 = this.selectionBound.y - this.containerPageTop
    let x2 = x1 + this.selectionBound.width
    let y2 = y1 + this.selectionBound.height

    const x1TimeBin = this.getTimeBinX(x1, true)
    const x2TimeBin = this.getTimeBinX(x2, false)
    const y1TimeBin = this.getTimeBinY(y1, false)
    const y2TimeBin = this.getTimeBinY(y2, true)

    return {
      x1TimeBin,
      x2TimeBin,
      y1TimeBin,
      y2TimeBin
    }
  }

  get animatedStyle () {
    const x1TimeBin = this.brushedTimeBin.x1TimeBin
    const x2TimeBin = this.brushedTimeBin.x2TimeBin
    const y1TimeBin = this.brushedTimeBin.y1TimeBin
    const y2TimeBin = this.brushedTimeBin.y2TimeBin

    let x1Refined = x1TimeBin / this.xTimeBin * this.timelineContainerWidth
    let x2Refined = x2TimeBin / this.xTimeBin * this.timelineContainerWidth
    let widthRefined = x2Refined - x1Refined
    let y1Refined = y1TimeBin / this.yTimeBin * this.timelineContainerHeight
    let y2Refined = y2TimeBin / this.yTimeBin * this.timelineContainerHeight
    let heightRefined = y2Refined - y1Refined

    return {
      left: (x1Refined + this.containerPageLeft) + 'px',
      top: (y1Refined + this.containerPageTop) + 'px',
      width: widthRefined + 'px',
      height: heightRefined + 'px',
      border: `1px solid ${
        Color('#54a0ff').alpha(0.7).darken(0.1)}`,
      'background-color': Color('#54a0ff').alpha(0.3).toString(),
      'border-radius': 4 + 'px'
    }
  }

  private retreiveBrushedTimeWindows () {
    const x1TimeBin = this.brushedTimeBin.x1TimeBin
    const x2TimeBin = this.brushedTimeBin.x2TimeBin
    const y1TimeBin = this.brushedTimeBin.y1TimeBin
    const y2TimeBin = this.brushedTimeBin.y2TimeBin

    const brushedTreeIndices: {[index: number]: boolean} = {}
    const timeWindows: ITimeWindowInfo[] = []
    for (let x = x1TimeBin; x < x2TimeBin; x++) {
      if (y2TimeBin >= y1TimeBin) { //  for each time windows, get valid tree
        const lowerBound = x * this.yTimeBin + y1TimeBin
        const upperBound = x * this.yTimeBin + y2TimeBin
        this.trees.forEach((t: PropagationTree) => {
          const rowOk: boolean = (this.selectedTreesIdDictByATrees[t.id] && this.selectedAggregatePropagationTreesId.length > 0) ||
            (this.selectedAggregatePropagationTreesId.length === 0)
          const columnOk: boolean = t.id in this.filteredTreesIndices
          if (rowOk && columnOk) {
            const treeStartTime = t.TreeStartTime
            if (treeStartTime > lowerBound && treeStartTime < upperBound) {
              brushedTreeIndices[t.id] = true
              timeWindows.push({ timeRange: [treeStartTime, t.TreeEndTime], tree: t })
            }
          }
        })
      }
    }
    if (timeWindows.length === 0) {
      this.resetSelection()
      this.resetTimeWindows(this.trees)
    } else {
      this.dispatchTimeWindows(timeWindows)
    }
  }

  get xTimeBin () {
    return this.currentDataSetConfig.xTimeBin
  }

  get yTimeBin () {
    return this.currentDataSetConfig.yTimeBin
  }

  get xTicks () {
    return new Array(this.currentDataSetConfig.gridShape[0]).fill(0).map((v, i) => i)
  }

  get yTicks () {
    return new Array(this.currentDataSetConfig.gridShape[1]).fill(0)
      .map((v, i) => this.currentDataSetConfig.yTimeBinRate * i)
  }

  get xScale (): (x: number) => number {
    return (x) => x * this.timelineContainerWidth / this.xTimeBin
  }

  get yScale (): (y: number) => number {
    return (y) => y * this.timelineContainerHeight / this.yTimeBin
  }

  private getTimeBinX (xpx: number, leftBottom: boolean) {
    const tx = xpx * this.xTimeBin / this.timelineContainerWidth
    if (leftBottom) {
      return Math.floor(tx)
    } else {
      return Math.ceil(tx)
    }
  }

  private getTimeBinY (ypx: number, leftBottom: boolean) {
    const ty = ypx * this.yTimeBin / this.timelineContainerHeight
    if (leftBottom) {
      return Math.floor(ty)
    } else {
      return Math.ceil(ty)
    }
  }

  private getPositionAlongX (treeStartTime: number) {
    return this.xScale(Math.floor(treeStartTime / this.yTimeBin))
  }

  private getPositionAlongY (treeStartTime: number) {
    return this.yScale(treeStartTime % this.yTimeBin)
  }

  private created () {
    this.columnsSum = new Array(this.xTicks.length + 1).fill(0)
    this.rowsSum = new Array(this.yTicks.length + 1).fill(0)
  }

  private mounted () {
    const rect = (this.$refs.timelineRef as HTMLElement).getBoundingClientRect()
    this.timelineContainerWidth = rect.width
    this.timelineContainerHeight = rect.height
    this.containerPageLeft = rect.left
    this.containerPageTop = rect.top
    const svg: d3.Selection<Element, unknown, any, any> = d3.select(this.$refs.treesSvgRef as Element)
    this.patternG = svg.append('g').attr('class', 'draw-g')
  }

  @Watch('filteredTreesIndices', { deep: true })
  private onFilterChanged (filteredTreesIndices: {[index: number]: boolean}) {
    this.patternG.selectAll('*').remove()
    this.drawRect(this.trees)
    this.resetTimeWindows(this.trees)
  }

  private resetTimeWindows (trees: PropagationTree[]): void {
    const timeWindows: ITimeWindowInfo[] = []
    trees.forEach((t) => {
      const rowOk: boolean = (this.selectedTreesIdDictByATrees[t.id] && this.selectedAggregatePropagationTreesId.length > 0) ||
        (this.selectedAggregatePropagationTreesId.length === 0)
      const columnOk: boolean = t.id in this.filteredTreesIndices
      if (rowOk && columnOk) {
        timeWindows.push({ timeRange: [t.TreeStartTime, t.TreeEndTime], tree: t })
      }
    })
    this.dispatchTimeWindows(timeWindows)
    this.resetSelection()
  }

  private drawRect (trees: PropagationTree[]): void {
    let headPadding: number = 3
    let rectPadding: number = 6
    if (this.currentDataSetConfig.gridShape[0] > 100) {
      headPadding = 1
      rectPadding = 3
    }
    this.initRowColumnSum()
    const width = this.timelineContainerWidth / this.xTimeBin - 2 * rectPadding

    trees.forEach((t: PropagationTree) => {
      const startX: number = Math.floor(t.TreeStartTime / this.yTimeBin)
      const endX: number = Math.floor(t.TreeEndTime / this.yTimeBin)
      const startY: number = this.getPositionAlongY(t.TreeStartTime)
      const endY: number = this.getPositionAlongY(t.TreeEndTime)
      const rowOk: boolean = (this.selectedTreesIdDictByATrees[t.id] && this.selectedAggregatePropagationTreesId.length > 0) ||
        (this.selectedAggregatePropagationTreesId.length === 0)
      const columnOk: boolean = t.id in this.filteredTreesIndices

      const opacity = rowOk && columnOk ? 1 : 0.1

      if (startX === endX) {
        this.patternG.append('rect')
          .attr('class', 'tree-rect')
          .attr('transform', `translate(${this.getPositionAlongX(t.TreeStartTime) + rectPadding},${this.getPositionAlongY(t.TreeStartTime)})`)
          .attr('fill', 'rgba(81, 97, 114, 0.3)')
          .attr('stroke-width', 1)
          .attr('stroke', 'rgba(87,101,116,0.7)')
          .attr('opacity', opacity)
          .attr('height', Math.max(endY - startY, 0))
          .attr('width', width)
      } else {
        this.patternG.append('rect')
          .attr('class', 'tree-rect')
          .attr('transform', `translate(${this.getPositionAlongX(t.TreeStartTime) + rectPadding},${this.getPositionAlongY(t.TreeStartTime)})`)
          .attr('fill', 'rgba(81, 97, 114, 0.3)')
          .attr('stroke-width', 1)
          .attr('stroke', 'rgba(87,101,116,0.7)')
          .attr('opacity', opacity)
          .attr('height', Math.max(this.timelineContainerHeight - startY, 0))
          .attr('width', width)
        let nowX = startX + 1
        while (nowX < endX) {
          this.patternG.append('rect')
            .attr('class', 'tree-rect')
            .attr('transform', `translate(${this.getPositionAlongX(t.TreeStartTime) + rectPadding + this.xScale(nowX - startX)},${0})`)
            .attr('fill', 'rgba(81, 97, 114, 0.3)')
            .attr('stroke-width', 1)
            .attr('stroke', 'rgba(87,101,116,0.7)')
            .attr('opacity', opacity)
            .attr('height', endY)
            .attr('width', width)
          nowX++
        }
        this.patternG.append('rect')
          .attr('class', 'tree-rect')
          .attr('transform', `translate(${this.getPositionAlongX(t.TreeStartTime) + rectPadding + this.xScale(1)},${0})`)
          .attr('fill', 'rgba(81, 97, 114, 0.3)')
          .attr('stroke-width', 1)
          .attr('stroke', 'rgba(87,101,116,0.7)')
          .attr('opacity', opacity)
          .attr('height', endY)
          .attr('width', width)
      }

      this.patternG.append('line')
        .attr('class', 'tree-head')
        .attr('transform', `translate(${this.getPositionAlongX(t.TreeStartTime)},${this.getPositionAlongY(t.TreeStartTime)})`)
        .attr('opacity', opacity)
        .attr('x1', headPadding)
        .attr('y1', 0)
        .attr('x2', this.timelineContainerWidth / this.xTimeBin - headPadding)
        .attr('y2', 0)
        .attr('stroke', 'rgba(87,101,116,0.7)')
        .attr('stroke-width', 2)
    })
  }

  @Watch('trees')
  private onTreesChanged (trees: PropagationTree[]) {
    this.patternG.selectAll('*').remove()
    this.drawRect(trees)
    this.resetTimeWindows(trees)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.dynamic-trees-view {
  position: relative;
  width: 100%;
  height: 100%;
  user-select: none;

  .title-div {
    height: 30px;
    width: 100%;
    background-color: white;
    text-align: left;
    align-items: center;

    .name {
      text-align: center;
      background-color: #0879eb;
      position: relative;
      top: 0px;
      left: 0px;
      width: 160px;
      height: 100%;
      font-size: 14px;
      font-family: Helvetica-Bold,Helvetica;
      font-weight: bold;
      color: white;
      line-height: 30px;
      letter-spacing: 1px;
      box-shadow: 0px 2px 4px 0px rgba(32, 84, 159, 0.4);
    }
  }

  .dynamic-trees-view-cont {
    position: relative;
    padding: 2px 20px 0px;
    box-sizing: border-box;
    height: calc(100% - 30px);
    width: 100%;
    .top-div {
      width: calc(100% - 22px);
      height: 20px;
      box-sizing: border-box;
      margin-bottom: 2px;
      .head-svg {
        width: 100%;
        height: 100%;
        background-color: #fff;
        .cell {
          height: 20px;
          fill: rgba(81, 97, 114, 1);
          stroke-width: 1;
          stroke: white;
        }
      }
    }

    .mid-div {
      width: 100%;
      height: calc(100% - 22px);
      background: #fff;
      display: flex;
      .trees-timeline {
        height: 100%;
        width: calc(100% - 20px);
        .timeline {
          position: relative;
          width: 100%;
          height: 100%;
          box-sizing: border-box;
          border: 1px solid #ddd;

          svg {
            position: relative;
            width: 100%;
            height: 100%;
          }
        }
      }
      .right-div {
        width: 20px;
        height: 100%;
        margin-left: 2px;
        .right-svg {
          width: 100%;
          height: 100%;
          background-color: #fff;
          .cell {
            width: 20px;
            fill: rgba(81, 97, 114, 1);
            stroke-width: 1;
            stroke: white;
          }
        }
      }
    }

    .brusher {
      position: absolute;
      left: 0px;
      top: 0px;
      width: 100%;
      height: 100%;
    }

    .selection {
      position: fixed;
    }

    .animated {
      transition: all 300ms;
    }
  }
}
</style>
