<template>
  <div class="projection">
    <div class="projection-canvas" ref="projCanvas">
      <svg :width="svgWidth" :height="svgHeight" class="proj-svg" ref="projSvg">
        <circle v-for="p in embedPoints" :key="p.rid" :id="p.rid"
                :cx="p.left * svgWidth" :cy="p.top * svgHeight"
                :r="3" :fill="'#777'" :opacity="getNormalizedValue(p.rid)"/>
      </svg>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'

import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import * as d3 from 'd3'
import { uniq, forEach } from 'lodash'
// @ts-ignore
import { lasso } from 'd3-lasso'
import EventBus from '@/eventBus'
import '@/css/lasso.css'

const projStore: BindingHelpers = namespace('Proj')

@Component
export default class TestProjectionInPresentation extends Vue {
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.State private allIDs!: number[]
  @projStore.Action private setProjectionCanvasSize!: (payload: {width: number, height: number}) => void
  @projStore.Action private getCorrelatedLocations!: (params: { clusterId: number, locationIDs: number[]}) => void
  @projStore.State private allLocationsInfoDict!: {[key: number]: ILocationInfo}
  @projStore.State private embedPoints!: ProjPoint[]
  @projStore.State private clustersWithContext!: ICluster[]
  @projStore.State private colorMapForEachLocation!: {[rid: number]: string}
  @projStore.Getter private filteredIDDict!: {[key: number]: number}
  private svgWidth: number = 420
  private svgHeight: number = 420
  private maxTotalCount: number = 1
  private minTotalCount: number = 0

  @Watch('allLocationsInfoDict')
  private setMinMaxTotalCount (): void {
    const totalCounts: number[] = this.allIDs.map((id: number) => {
      return this.allLocationsInfoDict[id].totalEventCount
    })
    this.maxTotalCount = Math.max(...totalCounts)
    this.minTotalCount = Math.min(...totalCounts)
    this.initLasso()
  }

  private selectLocations (locationsId: number[], clusterId: number): void {
    EventBus.$emit('onProjSelect', locationsId)
  }

  private mounted (): void {
    // const rect = (this.$refs.projCanvas as HTMLElement).getBoundingClientRect()
    // this.svgHeight = rect.height
    // this.svgWidth = rect.width
    this.setProjectionCanvasSize({ width: this.svgWidth, height: this.svgHeight })
    this.initLasso()
  }

  private initLasso (): void {
    const projSvg: d3.Selection<SVGElement, unknown, null, undefined> = d3.select(this.$refs.projSvg as SVGElement)
    const circles: d3.Selection<d3.BaseType, unknown, SVGElement, unknown> = projSvg.selectAll('circle')
    // Lasso functions
    const lassoStart = () => {
      myLasso.items().attr('r', 3)
      myLasso.items().attr('stroke', 'none')
    }

    // const lassoDraw = () => {
    //   myLasso.possibleItems().attr('fill', 'orange')
    //   myLasso.notPossibleItems().attr('fill', '#777')
    // }

    const lassoEnd = () => {
      const locationIds: number[] = []
      myLasso.selectedItems().nodes().forEach((v: any) => {
        if (typeof v === 'object') {
          locationIds.push(Number.parseInt(v.id))
        }
      })
      console.log('sampling: ', locationIds)
      if (uniq(locationIds).length > 3) {
        this.selectLocations(uniq(locationIds), 9999)
      }
      // Style the selected dots
      myLasso.selectedItems().attr('r', 6)
      myLasso.selectedItems().attr('stroke', '#fff')

      // Reset the style of the not selected dots
      myLasso.notSelectedItems().attr('r', 2)
      myLasso.notSelectedItems().attr('stroke', 'none')
    }

    const myLasso = lasso()
      .closePathSelect(true)
      .closePathDistance(100)
      .items(circles)
      .targetArea(projSvg)
      .on('start', lassoStart)
      // .on('draw', lassoDraw)
      .on('end', lassoEnd)

    projSvg.call(myLasso)
  }

  private getNormalizedValue (thisId: number): number {
    if (!(thisId in this.allLocationsInfoDict)) {
      return 0
    }
    const val: number = (this.allLocationsInfoDict[thisId].totalEventCount - this.minTotalCount) /
      (this.maxTotalCount - this.minTotalCount)
    return 0.1 + 0.9 * val
  }
}
</script>

<style lang="scss" scoped>
.projection {
  position: relative;
  width: 100%;
  height: 100%;
  border-right: 1px solid #eee;

  .projection-canvas {
    position: relative;
    width: 100%;
    height: 100%;
    border-bottom: none;
  }
}
</style>
