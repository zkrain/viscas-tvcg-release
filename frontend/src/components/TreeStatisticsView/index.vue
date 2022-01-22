<template>
  <div class="tree-statistics-view" ref="treeStatisticsView"  @scroll="onScroll($event)">
<!--    <graph-window class="graph-window"  ref="graphWindow" :edges-locations-id="edgesLocationsId" v-if="showWindow && clickedRight"-->
<!--                  :style="{left: (leftWindow + edgeWindowLeft) + 'px', top: (edgeWindowTop + topWindow) + 'px'}"/>-->
    <div class="head-container">
      <div class="view-title">
        <div class="name" :style="{left: (leftWindow + 20) + 'px'}">Influence View</div>
      </div>
      <div class="edge-information">
        <div class="edge-information-cont" v-for="(edge, ej) in cascadingLinks" :key="ej"
             @mouseenter="setHoveredEdge(`${edge[0]},${edge[1]}`)"
             @mouseleave="setHoveredEdge('')"
             @click="onSelectEdge(`${edge[0]},${edge[1]}`)"
             :style="{'background-color': selectedEdge[`${edge[0]},${edge[1]}`]? 'rgba(45, 52, 54, 0.13)': '#fff'}">

          <div class="edge-cont">
            <div class="location-i" :style="{ 'background-color': locationColorDict[edge[0]] }" />

            <div class="link-cont" :style="{
              width: `calc(${edgeLengthDict[`${edge[0]},${edge[1]}`] * 85}% - ${44 * edgeLengthDict[`${edge[0]},${edge[1]}`]}px)`
            }">
              <div class="link" :style="{
                  height: linkWidthRange[0] + (cascadingLinks.length - edge[2]) * (linkWidthRange[1] - linkWidthRange[0]) / cascadingLinks.length + 'px'
                }" /></div>

            <div class="location-j" :style="{ 'background-color': locationColorDict[edge[1]] }" />
            <transition name="fade">
              <div :class="[{ 'open-check': `${edge[0]},${edge[1]}` in openedEdge }, 'check-cont']"
                   v-if="hoveredEdgeKey === `${edge[0]},${edge[1]}` || `${edge[0]},${edge[1]}` in openedEdge"
                   @click.stop="onOpenEdge(`${edge[0]},${edge[1]}`)">
                <font-awesome-icon :icon="faCheckIcon" />
              </div>
            </transition>
          </div>
          <div class="p-distribution-cont">
            <div class="p-low number">{{pLowerBoundStr}}</div>
            <div class="p-high number">1</div>
            <div class="distribution-cont">
              <div v-for="(cell, cellIndex) in edgeProbabilityDistribution[`${edge[0]},${edge[1]}`]"
                   class="distribution-cell" :key="cellIndex"
                   :style="{ opacity: cell, 'background-color': locationColorDict[edge[0]] }"/>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="scroller-container" @click.right="clickedRight = !clickedRight">
      <div class="aggregate-trees-container">
        <transition-group name="list-complete" tag="div" mode="out-in">
          <div class="aggregate-tree" v-for="atree in aggregatePropagationTreesFilter" :key="atree.id"
               :class="{'aggregate-tree-select': atree.selected}"
               :style="{ height: `${atree.height}px`}"
               @click="selectATrees(atree)"
               @mousemove="onMouseMoveLine(atree, $event)"
               @mouseenter="onMouseEnterLine(atree, $event)"
               @mouseleave="showWindow = false">
            <div class="aggregate-tree-edge" v-for="(edge, j) in cascadingLinks" :key="j" :id="`${edge[0]},${edge[1]}`">
              <edge-cell :edge="edge" :atree="atree" />
            </div>
          </div>
        </transition-group>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faCheck } from '@fortawesome/free-solid-svg-icons'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import EdgeCell from './EdgeCell.vue'
import GraphWindow from '@/components/TreeStatisticsView/GraphWindow.vue'

const projStore: BindingHelpers = namespace('Proj')

@Component({
  components: {
    GraphWindow,
    FontAwesomeIcon,
    EdgeCell
  }
})
export default class TreeStatisticsView extends Vue {
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.State private url!: string
  @projStore.State private cascadingLinks!: number[][]
  @projStore.State private aggregatePropagationTrees!: AggregatePropagationTree[]
  @projStore.State private edgeProbabilityDistribution!: {[edgeKey: string]: number[]}
  @projStore.State private openedEdge!: {[edgeKey: string]: boolean}
  @projStore.State private selectedEdge!: {[edgeKey: string]: boolean}
  @projStore.State private selectedAggregatePropagationTreesId!: number[]
  @projStore.Getter private selectedTreesIdDictByATrees!: {[key: number]: boolean}
  @projStore.Getter private filteredTreesIndices!: {[index: number]: boolean}
  @projStore.Getter private edgeColorDict!: {[edgeKey: string]: [string, string]}
  @projStore.Getter private locationColorDict!: {[locationID: number]: string}
  @projStore.Getter private edgeLengthDict!: {[key: string]: number}
  @projStore.Action private selectEdge!: (edgeKey: string) => void
  @projStore.Action private openEdge!: (edgeKey: string) => void
  @projStore.Action private selectTrees!: (trees: PropagationTree[]) => void
  @projStore.Action private selectAtree!: (ind: number) => void
  @projStore.State private cascadingPointsScreenPosition!: {[locationID: number]: [number, number]}

  private aggregatePropagationTreesFilter: AggregatePropagationTree[] = []
  private linkWidthRange: [number, number] = [4, 8]
  private hoveredEdgeKey: string = ''

  private edgesLocationsId: [number, number][] = []
  private edgeWindowLeft: number = 0
  private edgeWindowTop: number = 0
  private domLeft: number = 0
  private domTop: number = 0
  private domHeight: number = 0
  private domWidth: number = 0
  private windowHeight: number = 140
  private windowWidth: number = 140
  private showWindow: boolean = false
  private clickedRight: boolean = true
  private topWindow: number = 0
  private leftWindow: number = 0

  get pLowerBoundStr () {
    return this.currentDataSetConfig.pLowerBound.toString().slice(1)
  }

  get faCheckIcon () {
    return faCheck
  }

  private mounted (): void {
    const ref = this.$refs.treeStatisticsView as HTMLElement
    const rect = ref.getBoundingClientRect()
    this.domLeft = rect.left
    this.domTop = rect.top
    this.domHeight = rect.height
    this.domWidth = rect.width
    ref.oncontextmenu = function (e) {
      e.preventDefault()
    }
  }

  private onScroll (e: Event): void {
    const target = e.currentTarget as HTMLElement
    this.topWindow = target.scrollTop
    this.leftWindow = target.scrollLeft
  }

  private onMouseEnterLine (atree: AggregatePropagationTree, e: MouseEvent): void {
    this.edgesLocationsId = []
    atree.involvedEdges.forEach((edge) => {
      this.edgesLocationsId.push([edge[0], edge[1]])
    })
    this.showWindow = true
  }

  private onMouseMoveLine (atree: AggregatePropagationTree, e: MouseEvent): void {
    this.edgeWindowLeft = e.pageX - this.domLeft + 20
    this.edgeWindowTop = e.pageY - this.domTop + 20
    if (this.edgeWindowLeft + 20 + this.windowWidth > this.domWidth) {
      this.edgeWindowLeft = this.edgeWindowLeft - 40 - this.windowWidth
    }
    if (this.edgeWindowTop + 20 + this.windowHeight > this.domHeight) {
      this.edgeWindowTop = this.edgeWindowTop - 40 - this.windowHeight
    }
  }

  @Watch('aggregatePropagationTrees')
  @Watch('selectedEdge', { deep: true })
  private async setAggregatePropagationTreesFilter (): Promise<void> {
    if (typeof this.aggregatePropagationTrees === 'undefined') {
      this.aggregatePropagationTreesFilter = []
      return
    }
    let tAggregatePropagationTreesFilter: AggregatePropagationTree[] = []
    let selectedEdgeNum: number = 0
    for (const key in this.selectedEdge) {
      if (this.selectedEdge[key]) {
        selectedEdgeNum++
      }
    }
    if (selectedEdgeNum === 0) {
      tAggregatePropagationTreesFilter = this.aggregatePropagationTrees
    } else {
      tAggregatePropagationTreesFilter = this.aggregatePropagationTrees.filter((t: AggregatePropagationTree) => {
        let hasEdgeNum: number = 0
        for (const key in this.selectedEdge) {
          if (key in t.involvedEdgeDictInAggregation && this.selectedEdge[key]) {
            hasEdgeNum++
          }
        }
        return selectedEdgeNum === hasEdgeNum
      })
    }

    const vertices: number[][] = []
    const heights: number[] = []
    tAggregatePropagationTreesFilter.forEach((atree: AggregatePropagationTree) => {
      const v: number[] = []
      this.cascadingLinks.forEach((edge: number[], j: number) => {
        const key: string = `${edge[0]},${edge[1]}`
        if (key in atree.involvedEdgeDictInAggregation) {
          v.push(j)
        }
      })
      vertices.push(v)
      heights.push(atree.height)
    })

    const respHamiltonian = await Vue.http.post(`${this.url}hamiltonianWalkOptimizer`, { vertices, heights })
    const order: number[] = respHamiltonian.data.order

    // 重新排序
    this.aggregatePropagationTreesFilter = []
    order.forEach((val: number) => {
      this.aggregatePropagationTreesFilter.push(tAggregatePropagationTreesFilter[val])
    })
  }

  private selectATrees (atree: AggregatePropagationTree): void {
    this.selectAtree(atree.id)
    this.selectTrees(atree.trees)
  }

  private onSelectEdge (edgeKey: string) {
    this.selectEdge(edgeKey)
  }

  private onOpenEdge (edgeKey: string) {
    this.openEdge(edgeKey)
  }

  private setHoveredEdge (edgeKey: string) {
    this.hoveredEdgeKey = edgeKey
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
$EDGEWIDTH: 147px;
// $EDGEWIDTH: 130px;

.tree-statistics-view {
  position: relative;
  width: calc(100% - 3px);
  height: calc(100% - 2px);
  margin: 1px;
  border-right: 1px solid #eee;
  overflow-x: scroll;
  overflow-y: hidden;

  .graph-window {
    position: absolute;
    z-index: 999;
  }

  .head-container {
    position: relative;
    width: 1300px;
    .view-title {
      width: 100%;
      height: 30px;
      position: relative;
      background:rgba(8,121,235,1);
      box-shadow:0px 2px 4px 0px rgba(32,84,159,0.4);

      .name {
        position: absolute;
        top: 6px;
        font-size:14px;
        font-family:Helvetica-Bold,Helvetica;
        font-weight:bold;
        color:rgba(255,255,255,1);
      }
    }
    .edge-information {
      position: relative;
      height: 39px;
      width: 1300px;
      border-bottom: 1px solid #ddd;

      .edge-information-cont {
        position: relative;
        height: calc(100% - 3px);
        width: calc(#{$EDGEWIDTH} - 4px);
        float: left;
        background-color: #fff;
        margin: 2px 2px 0 2px;
        border: 1px solid #c9c9c9;
        border-bottom: none;
        border-top-left-radius: 4px;
        border-top-right-radius: 4px;

        .edge-cont {
          position: relative;
          height: calc(100% - 14px);
          width: calc(100% - 12px);
          margin: 2px 6px;

          .location-i {
            position: relative;
            height: 20px;
            width: 20px;
            border-radius: 50%;
            float: left;
            border: 1px solid #fff;
          }

          .location-j {
            position: relative;
            float: left;
            height: 20px;
            width: 20px;
            border-radius: 50%;
            border: 1px solid #fff;
          }

          .link-cont {
            float: left;
            position: relative;
            height: 100%;

            .link {
              position: relative;
              width: 100%;
              top: 50%;
              transform: translate(0, -50%);
              background-color: #576574;
              border-top: 1px solid #fff;
              border-bottom: 1px solid #fff;
            }
          }

          .check-cont {
            position: absolute;
            top: -3px;
            left: calc(100% - 12px);
            width: 16px;
            height: 16px;
            color: #bbb;
            transition: color 300ms;
          }

          .check-cont:hover {
            color: #00a8ff;
          }

          .open-check {
            color: #00a8ff !important;
          }

        }

        .p-distribution-cont {
          position: relative;
          height: 10px;
          width: 100%;

          .p-low {
            left: 0px;
          }

          .p-high {
            left: calc(100% - 6px);
          }

          .number {
            position: absolute;
            top: -3px;
            color: #576574;
            font-size: 12px;
            transform: scale(0.8);
          }

          .distribution-cont {
            position: relative;
            width: 100%;
            height: 100%;
            display: flex;

            .distribution-cell {
              flex: 1;
              background-color: #576574;
            }
          }
        }
      }
    }
  }

  .scroller-container {
    position: relative;
    width: 1300px;
    height: calc(100% - 42px);
    overflow-y: scroll;

    .aggregate-trees-container {
      position: relative;
      height: 100%;
      /*overflow-y: scroll;*/

      .aggregate-tree {
        position: relative;
        vertical-align: middle;
        border-bottom: 1px solid #eee;
        margin-bottom: 2px;
        transition: all 500ms;

        &:hover {
          border-color: rgb(28, 197, 240);
        }

        .aggregate-tree-edge {
          position: relative;
          float: left;
          width: $EDGEWIDTH;
          background-color: #fff;
          height: calc(100% - 2px);
          margin: 0 1px;
        }
      }
      .aggregate-tree-select {
        border-color: rgb(28, 197, 240);
      }
    }

  }
}

.list-complete-enter, .list-complete-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
.list-complete-leave-active {
  position: relative;
}

::-webkit-scrollbar {
  display: none;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity .3s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
</style>
