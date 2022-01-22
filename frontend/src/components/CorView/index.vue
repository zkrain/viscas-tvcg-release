<template>
  <div class="corview">
    <div class="corview-menu">
      <div class="reduce-noise-btn corview-menu-btn" @click="onNoiseReduce">
        Reduce Noise
      </div>
      <div class="cascading-infer-btn corview-menu-btn" @click="onCascadingInfer">
        Inferring
      </div>
    </div>
    <!-- <div class="infer-btn" @click="onInfer">
      Infer
    </div> -->

    <!-- <div class="layout-btn" @click="onReLayout">
      Relayout
    </div> -->
    <div class="corview-cont">
      <div ref="event-background-cont" class="event-background">
        <div v-for="l in timeStamps" :key="l" class="vert-line"/>
      </div>
      <div class="events-with-id-list">
        <div v-for="d in eventsListSlice" :key="d.locationID" class="events-with-id">
          <!-- <div class="events" v-for="d in eventsListSlice" :key="d.locationID"
          @mouseenter="enterRoadRow(r.id)"
          @mouseleave="leaveRoadRow"> -->
          <div class="locationID">
            {{ d.locationID }}
          </div>
          <div class="events">
            <div class="slot" />
            <div v-for="e in d.events" :key="e" class="event" :title="`${e}_${d.locationID}`"
                 :style="{ left: `${e*widthPerTime - 2}px` }"/>
          </div>
        </div>
        <div ref="svg-cont" class="svg-cont" :style="{ height: `${svgHeight}px` }">
          <svg :width="svgWidth" :height="svgHeight">
            <g v-for="(InfluencesWithID, i) in InfluencesWithIDListSlice" :key="i">
              <line v-for="(e, j) in InfluencesWithID.Influences" :key="j"
                    :style="{ strokeWidth: e.P * 1, stroke: 'red' }" :x1="widthPerTime * e.StartTimeI"
                    :y1="heightPerRow * IDIndex[InfluencesWithID.locationIDi] + heightPerRow / 2"
                    :x2="widthPerTime * e.StartTimeJ"/>
            </g>
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { Action, Getter, namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import EventBus from '../../eventBus'
import { sortBy, filter, forEach } from 'lodash'

const projStore: BindingHelpers = namespace('Proj')

@Component
export default class CorView extends Vue {
  @projStore.State private eventsWithIDList!: IEventsWithID[]
  // @projStore.Action private inferCascadingPatternsPatterns!: () => void
  // @projStore.Action private reLayout!: () => void
  @projStore.Action private reduceNoiseForEvents!: (payload: { window: number, threshold: number }) => void
  // @projStore.State private inspectingRoads!: Road[]
  // @projStore.State private edges!: Edges
  @projStore.State private timeLength!: number
  @projStore.Getter private InfluencesWithIDListSlice!: InfluencesWithID[]

  private eventContWidth: number = 1
  private svgWidth: number = 1
  private heightPerRow: number = 18

  // // private edgesList: InfluencesWithID[] = []

  get widthPerTime (): number {
    return this.eventContWidth / this.timeLength
  }

  get IDIndex (): {[roadID: number]: number} {
    const index: {[roadID: number]: number} = {}
    this.eventsWithIDList.forEach((eventsWithID, i) => {
      index[eventsWithID.locationID] = i
    })
    return index
  }

  get timeStamps (): number[] {
    return [...Array(this.timeLength).keys()]
  }

  get eventsListSlice (): IEventsWithID[] {
    return this.eventsWithIDList.map((d: IEventsWithID) => {
      const eventsSlice: number[] = filter(d.events, (e: number) => e < this.timeLength)
      return {
        events: eventsSlice,
        locationID: d.locationID
      }
    })
  }

  get svgHeight (): number {
    return this.eventsWithIDList.length * this.heightPerRow
  }

  // private onReLayout() {
  //   this.reLayout()
  // }

  // private enterRoadRow (roadID: number) {
  //   EventBus.$emit('hightlightRoad', roadID)
  // }

  // private leaveRoadRow () {
  //   EventBus.$emit('cancelHighlightRoad')
  // }

  private async onNoiseReduce () {
    await this.reduceNoiseForEvents({
      window: 5,
      threshold: 0.4
    })
  }

  private async onCascadingInfer () {
    // this.inferCascadingPatternsPatterns()
  }

  private mounted () {
    const rectEventCont = (this.$refs['event-background-cont'] as HTMLElement).getBoundingClientRect()
    this.eventContWidth = rectEventCont.width

    const rectSvgCont = (this.$refs['svg-cont'] as HTMLElement).getBoundingClientRect()
    this.svgWidth = rectSvgCont.width

    // EventBus.$on('generateEventHeatmap', (congestEvents: number[][]) => {
    //   const canvas = document.createElement('canvas')
    //   const ctx = <CanvasRenderingContext2D>canvas.getContext('2d')

    //   canvas.height = congestEvents.length
    //   canvas.width = 8785
    //   ctx.fillStyle = '#FF0000'

    //   sortBy(congestEvents, (o) => o.length).forEach((times, index) => {
    //     times.forEach((t) => {
    //       ctx.fillRect(t, index, 1, 1)
    //     })
    //   })
    //   const canvasData = canvas.toDataURL('image/png', 1)
    // })
  }
}
</script>

<style lang="scss">
$id-width: 60px;
$row-height: 18px;

.corview {
  position: relative;
  height: calc(100% - 2px);
  width: calc(100% - 2px);
  border: 1px solid #eee;
  border-radius: 3px;

  // .infer-btn {
  //   position: absolute;
  //   height: 30px;
  //   width: 100px;
  //   left: 10px;
  //   top: -30px;
  //   background-color: #555;
  //   border-radius: 2px;
  //   color: #fff;
  //   line-height: 30px;
  //   cursor: pointer;
  // }
  .corview-menu {
    position: relative;
    width: 100%;
    height: 36px;

    .corview-menu-btn {
      position: relative;
      float: left;
      height: 30px;
      width: 150px;
      margin: 2px 10px;
      border: 1px solid #eee;
      background-color: #555;
      border-radius: 2px;
      color: #fff;
      line-height: 30px;
      cursor: pointer;
    }
  }

  // .layout-btn {
  //   position: absolute;
  //   height: 30px;
  //   width: 100px;
  //   left: 290px;
  //   top: -30px;
  //   background-color: #555;
  //   border-radius: 2px;
  //   color: #fff;
  //   line-height: 30px;
  //   cursor: pointer;
  // }

  .corview-cont {
    position: relative;
    height: calc(100% - 20px - 36px);
    width: calc(100% - 20px);
    margin: 10px;
    display: flex;

    .events-with-id-list {
      position: relative;
      height: 100%;
      width: 100%;
      overflow: scroll;

      .events-with-id {
        position: relative;
        height: $row-height;
        width: 100%;
        display: flex;

        .locationID {
          position: relative;
          flex: 0 0 $id-width;
          line-height: $row-height;
          text-align: start;
          border-right: 1px solid #f3f3f3;
          font-size: 14px;
        }

        .events {
          position: relative;
          flex: 1;

          .event {
            position: absolute;
            top:calc(50% - 1px);
            height: 2px;
            width: 2px;
            background-color: #d21e1e;
          }

          .slot {
            position: absolute;
            height: 6px;
            width: 100%;
            left: 0px;
            top: calc(50% - 3px);
            background-color: #eee;
          }
        }
      }

      .svg-cont {
        position: absolute;
        width: calc(100% - #{$id-width});
        left: $id-width;
        top: 0;
      }
    }

    .roads::-webkit-scrollbar {
      display: none;
    }

    .event-background {
      position: absolute;
      width: calc(100% - #{$id-width});
      height: 100%;
      left: $id-width;
      top: 0;
      background-color: #fafafa;
      display: flex;

      .vert-line {
        flex: 1;
        border-left: 0.5px solid #eee;
      }
    }
  }
}
</style>
