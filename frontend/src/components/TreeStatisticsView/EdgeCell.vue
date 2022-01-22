<template>
  <div class="edge-cell">
    <div class="aggregate-tree-edge-probability"
      v-if="`${edge[0]},${edge[1]}` in atree.involvedEdgeDictInAggregation">

      <div class="uncertainty content animated"
        :class="[{ closed: isAtreeSelected || isEdgeOpened }]"
        :style="{
          opacity: uncertaintyOpacity * 0.9 + 0.1
        }">
        <!-- <div class='left-horizon-line' v-for="(leftHorizonLineStyle, si) in leftHorizonLines" :key="`${si}_left-horizon-line`"
          :style="leftHorizonLineStyle"/>

        <div class="influence"
          v-for="(influenceInfoSpereated, ii) in influenceInfoSpereatedList"
          :key="ii">

          <div class='left5'>
            <div class="p-duration" :style="{
              width: `calc(${getLeft(influenceInfoSpereated.left.p) * 100}%)`,
              'border-color': locationColorDict[edge[0]]
            }" />
            <div class="opc" v-for="(opc, iiii) in influenceInfoSpereated.left.opcs" :key="iiii"
              :style="{
                left: `calc(${getLeft(opc.P) * 100}%)`,
                opacity: (opc.Duration + 1) >= influenceInfoSpereated.left.duration ? 0.8 : 0,
                backgroundColor: opc.LocationID in locationColorDict ? locationColorDict[opc.LocationID] : 'rgba(81, 97, 114, 0.6)'
              }" />
          </div>
        </div> -->
      </div>

      <div class="influences content animated"
        :class="[{ closed: !(isAtreeSelected || isEdgeOpened) }]">

        <div class="influence"
            v-for="(influenceInfoSpereated, ii) in influenceInfoSpereatedList"
            :key="ii">

          <div class='right2'>
            <div v-for="(durationLink, iiii) in influenceInfoSpereated.right" :key="iiii"
              :class="[{ overlap: durationLink[2] === -1 }, 'duration-link']"
              :style="{
                left: `calc(${100 * durationLink[0] / currentDataSetConfig.totalInfluenceDuration}%)`,
                width: `calc(${100 * (durationLink[1] - durationLink[0]) / currentDataSetConfig.totalInfluenceDuration}%)`,
                backgroundColor: durationLink[2] === -1
                  ? mixColor(locationColorDict[edge[1]], locationColorDict[edge[0]])
                  : locationColorDict[edge[durationLink[2]]]
              }"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch, Prop } from 'vue-property-decorator'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import { indexOf } from 'lodash'
import Color from 'color'

const projStore: BindingHelpers = namespace('Proj')

interface pDurationCause {
  p: number
  duration: number
  opcs: PotentialCause[]
}

interface InfluenceInfoSpereated {
  left: pDurationCause
  right: [number, number, number][]
}

@Component({
  components: {
    FontAwesomeIcon
  }
})
export default class EdgeCell extends Vue {
  @Prop() private atree!: AggregatePropagationTree
  @Prop() private edge!: [number, number]
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.State private selectedAggregatePropagationTreesId!: number[]
  @projStore.State private openedEdge!: {[edgeKey: string]: boolean}
  @projStore.Getter private locationColorDict!: {[locationID: number]: string}

  get isAtreeSelected () {
    return indexOf(this.selectedAggregatePropagationTreesId, this.atree.id) >= 0
  }

  get isEdgeOpened () {
    return `${this.edge[0]},${this.edge[1]}` in this.openedEdge
  }

  private mixColor (c1: string, c2: string): string {
    const color1 = Color(c1).object()
    const color2 = Color(c2).object()
    const r = (color1.r + color2.r) / 2
    const g = (color1.g + color2.g) / 2
    const b = (color1.b + color2.b) / 2
    return `rgba(${r},${g},${b},0.7)`
  }

  get uncertaintyOpacity () {
    const influenceInfoList = this.atree.involvedEdgeDictInAggregation[`${this.edge[0]},${this.edge[1]}`]
    const pDurationCauses = influenceInfoList
      .concat([])
      .map(info => ({ p: info.p, duration: info.duration, opcs: info.opcs }))
    const numberOfInfluences = influenceInfoList.length
    let numberOfOtherPossibilityWithLongDuration = 0
    for (const pDurationCause of pDurationCauses) {
      for (const opc of pDurationCause.opcs) {
        if (opc.Duration > pDurationCause.duration && opc.P > pDurationCause.p) {
          numberOfOtherPossibilityWithLongDuration += opc.P / pDurationCause.p
        }
      }
    }

    return numberOfOtherPossibilityWithLongDuration / numberOfInfluences
  }

  get influenceInfoSpereatedList () {
    const influenceInfoSpereatedList: InfluenceInfoSpereated[] = []
    const influenceInfoList = this.atree.involvedEdgeDictInAggregation[`${this.edge[0]},${this.edge[1]}`]
    const pDurationCauses = influenceInfoList
      .concat([])
      .sort((info1, info2) => info1.p - info2.p)
      .map(info => ({ p: info.p, duration: info.duration, opcs: info.opcs }))
    const influenceLinks = influenceInfoList.map(info => info.durationLinksAlignByRstj)
    for (let i = 0; i < influenceInfoList.length; i++) {
      influenceInfoSpereatedList.push({
        left: pDurationCauses[i],
        right: influenceLinks[i]
      })
    }

    return influenceInfoSpereatedList
  }

  get leftHorizonLines () {
    const leftHorizonLines = []
    for (let i = 0; i < this.influenceInfoSpereatedList.length - 1; i++) {
      const p1 = this.influenceInfoSpereatedList[i].left.p
      const p2 = this.influenceInfoSpereatedList[i + 1].left.p
      if (p1 !== p2) {
        leftHorizonLines.push({
          top: `${(i + 1) * 4}px`,
          width: `calc(${100 * this.getLeft(p2) - 100 * this.getLeft(p1)}%)`,
          left: `calc(${this.getLeft(p1) * 100}%)`,
          'background-color': this.locationColorDict[this.edge[0]]
        })
      }
    }
    return leftHorizonLines
  }

  private getLeft (p: number) {
    return (p - this.currentDataSetConfig.pLowerBound) / (1 - this.currentDataSetConfig.pLowerBound)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.edge-cell {
  position: relative;
  width: 100%;
  height: 100%;

  .aggregate-tree-edge-probability {
    position: absolute;
    top: 0px;
    left: 0px;
    height: calc(100% - 4px);
    width: calc(100% - 4px);
    border: 1px solid rgba(189, 195, 199, .7);
    // background-color: rgba(189, 195, 199, .1);
    margin: 1px;
    // box-shadow: 0px 0px 5px 1px rgba(199,199,199,1);

    .p-hint {
      position: relative;
      height: 100%;
      width: 0px;
      border-left: 1px solid rgba(87, 101, 116, 0.5);
    }

    .open {
      width: 100%;
    }

    .closed {
      // width: 0px !important;
      // border-right: none;
      opacity: 0 !important;
      overflow: hidden;
    }

    .animated {
      transition: all 300ms;
    }

    .content {
      position: absolute;
      top: 0px;
      left: 0px;
      width: 100%;
      height: 100%;
      display: flex;
      flex-direction: column;
      background-color: rgb(189, 195, 199);
    }

    .influences {
      background-color: rgba(189, 195, 199, .1) !important;
    }

    .left-horizon-line {
      position: absolute;
      height: 2px;
      background-color: #516172;
      transition: width 300ms;
    }

    .left-horizon-line.closed {
      width: 0px !important;
    }

    .influence {
      position: relative;
      flex: 1;
      height: 100%;
      width: 100%;
      overflow: hidden;

      .left5 {
        position: relative;
        width: calc(100%);
        height: 100%;
        border-right: 1px solid #ddd;
        float: left;

        .p-duration {
          position: relative;
          height: 100%;
          border-right: 2px solid #516172;
        }

        .opc {
          position: absolute;
          background-color: rgba(81, 97, 114, 0.6);
          width: 4px;
          height: 4px;
          top: 0px;
          border-radius: 50%;
        }
      }

      .right2 {
        position: relative;
        width: 100%;
        height: 100%;
        float: left;

        .duration-link {
          position: absolute;
          height: 1px;
          top: calc(50% - 0.5px);
          opacity: 1;
        }

        .overlap {
          height: 3px;
          top: calc(50% - 1.5px);
          opacity: 1;
        }
      }
    }
  }
}
</style>
