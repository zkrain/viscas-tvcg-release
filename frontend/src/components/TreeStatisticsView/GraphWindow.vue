<template>
  <div class="graph-container">
    <svg class="graph-svg" ref="graphSVG">
      <defs>
        <marker id="smallArrow" markerWidth="6" markerHeight="6" refX="8" refY="2" orient="auto" markerUnits="strokeWidth">
          <path d="M0,0 L0,4 L6,2 z" fill="#777" />
        </marker>
      </defs>
      <line v-for="(edge, key) in normalizedEdgesPosition" :key="`e-${key}`" stroke-width="2" stroke="#777"
            :x1="edge[0][0]" :y1="edge[0][1]" :x2="edge[1][0]" :y2="edge[1][1]" marker-end="url(#smallArrow)"/>
      <circle v-for="(pos, key) in normalizedLocationPosition" :key="`r-${key}`" r="6" :cx="pos[0]" :cy="pos[1]"
      :fill="locationColorDict[pos[2]]"/>
    </svg>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { flattenDeep, uniq } from 'lodash'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'

const projStore: BindingHelpers = namespace('Proj')

@Component
export default class GraphWindow extends Vue {
  @Prop() private edgesLocationsId!: [number, number][]
  @projStore.Getter private locationsScreenBound!: number[]
  @projStore.State private cascadingPointsScreenPosition!: {[locationID: number]: [number, number]}
  @projStore.Getter private locationColorDict!: {[locationID: number]: string}
  private rectWidth: number = 0
  private rectHeight: number = 0
  private mounted (): void {
    const rect = (this.$refs.graphSVG as HTMLElement).getBoundingClientRect()
    this.rectWidth = rect.width
    this.rectHeight = rect.height
  }
  get normalizedEdgesPosition (): [number, number][][] {
    if (typeof this.edgesLocationsId === 'undefined') {
      return []
    }
    const xRange: [number, number] = [this.locationsScreenBound[0] - 20, this.locationsScreenBound[1] + 20]
    const yRange: [number, number] = [this.locationsScreenBound[2] - 20, this.locationsScreenBound[3] + 20]
    const maxRangeLength: number = Math.max(yRange[1] - yRange[0], xRange[1] - xRange[0])
    const midX: number = 0.5 * (xRange[1] + xRange[0])
    const midY: number = 0.5 * (yRange[1] + yRange[0])
    return this.edgesLocationsId.map((locationsId: [number, number]) => {
      const edge: [number, number][] = [this.cascadingPointsScreenPosition[locationsId[0]],
        this.cascadingPointsScreenPosition[locationsId[1]]]
      return [[this.rectWidth * ((edge[0][0] - midX) / maxRangeLength + 0.5),
        this.rectHeight * ((edge[0][1] - midY) / maxRangeLength + 0.5)],
      [this.rectWidth * ((edge[1][0] - midX) / maxRangeLength + 0.5),
        this.rectHeight * ((edge[1][1] - midY) / maxRangeLength + 0.5)]]
    }) as [number, number][][]
  }

  get normalizedLocationPosition (): [number, number, number][] {
    if (typeof this.edgesLocationsId === 'undefined') {
      return []
    }
    const xRange: [number, number] = [this.locationsScreenBound[0] - 20, this.locationsScreenBound[1] + 20]
    const yRange: [number, number] = [this.locationsScreenBound[2] - 20, this.locationsScreenBound[3] + 20]
    const maxRangeLength: number = Math.max(yRange[1] - yRange[0], xRange[1] - xRange[0])
    const midX: number = 0.5 * (xRange[1] + xRange[0])
    const midY: number = 0.5 * (yRange[1] + yRange[0])
    const uniqLocationsId: number[] = uniq(flattenDeep(this.edgesLocationsId)) as number[]
    return uniqLocationsId.map((id) => {
      const pos: [number, number] = this.cascadingPointsScreenPosition[id]
      return [this.rectWidth * ((pos[0] - midX) / maxRangeLength + 0.5),
        this.rectHeight * ((pos[1] - midY) / maxRangeLength + 0.5), id]
    })
  }
}
</script>

<style scoped lang="scss">
.graph-container {
  width: 100px;
  height: 100px;
  border: 1px solid #ddd;
  background-color: rgba(255, 255, 255, 0.6);
  border-radius: 4px;
  pointer-events: none;
  .graph-svg {
    width: 100%;
    height: 100%;
  }
}
</style>
