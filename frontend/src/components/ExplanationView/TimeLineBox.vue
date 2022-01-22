<template>
  <div class="time-line-box-container"
       :style="{width: currentDataSetConfig.timelineLengthRate * (timeWindowInfo.timeRange[1] - timeWindowInfo.timeRange[0]) + 'px'}">
    <div class="head" :title="`${startTime}\n${endTime}`">
      <div class="time-txt">{{startTime}}</div>
      <div class="time-txt">{{endTime}}</div>
    </div>
<!--    <div  class="timeline-box">-->
<!--      <div v-for="(row, i) in timelines" :key="i" class="line">-->
<!--        <div v-for="(t, j) in row" :key="j" class="cell"-->
<!--             :class="{'first-cell': isFirst(j, row), 'last-cell': isLast(j, row)}"-->
<!--             :style="{'background-color': (t===1?'#777' : 'inherit')}"/>-->
<!--        &lt;!&ndash; todo: use one div &ndash;&gt;-->
<!--      </div>-->
<!--    </div>-->
    <div  class="timeline-box">
      <div v-for="(row, i) in infectEventsLines" :key="i" class="line">
        <div v-for="(t, j) in row" :key="j" class="cell"
             :class="{'cell-event': t[0]===1}"
             :style="{'flex': t[1],
             'height': t[0] === 2? '20%' : '100%',
             'background-color': t[0] >=1? locationColorDict[locationsId[i]]:'inherit',
             'opacity': t[4]}"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import { getDateFromNumber } from '@/utils/draw'
import { find } from 'lodash'

const projStore: BindingHelpers = namespace('Proj')

@Component
export default class TimeLineBox extends Vue {
  @Prop() private timeWindowInfo!: ITimeWindowInfo
  @Prop() private locationsId!: number[]
  @projStore.Getter locationColorDict!: {[locationID: number]: string}
  @projStore.State clusterInfectEventsDict!: {[key: number]: InfectEvent[]}
  @projStore.State currentDataSetConfig!: IDataSetConfig
  private infectEventsLines: number[][][] = [] // [infectEvent or not, duration, startTime, endTime, opacity]

  private mounted (): void {
    this.initTimelines()
  }

  get startTime (): string {
    return getDateFromNumber(this.currentDataSetConfig.startTime,
      this.currentDataSetConfig.deltaTime, this.timeWindowInfo.timeRange[0])
  }

  get endTime (): string {
    return getDateFromNumber(this.currentDataSetConfig.startTime,
      this.currentDataSetConfig.deltaTime, this.timeWindowInfo.timeRange[1])
  }

  // 这个文件里的注释很重要，不要删除任何一行
  // private isFirst (ind: number, line: number[]): boolean {
  //   if (line[ind] !== 1) {
  //     return false
  //   }
  //   if (ind === 0) {
  //     return true
  //   }
  //   return line[ind - 1] === 0
  // }
  //
  // private isLast (ind: number, line: number[]): boolean {
  //   if (line[ind] !== 1) {
  //     return false
  //   }
  //   if (ind === line.length - 1) {
  //     return true
  //   }
  //   return line[ind + 1] !== 1
  // }

  private initTimelines (): void {
    // this.timelines = []
    this.infectEventsLines = []
    const timeRange: [number, number] = this.timeWindowInfo.timeRange
    this.locationsId.forEach((val: number) => {
      const infectEvents: InfectEvent[] = this.clusterInfectEventsDict[val]
      const eLen: number = infectEvents.length
      // let line: number[] = fill(Array(this.timeRange[1] - this.timeRange[0] + 1), 0)
      // for (let i = 0; i < eLen; i++) {
      //   const e = infectEvents[i]
      //   if (e.startTime > this.timeRange[1]) {
      //     break
      //   }
      //   if (e.endTime < this.timeRange[0]) {
      //     continue
      //   }
      //   const start: number = Math.max(e.startTime - this.timeRange[0], 0)
      //   const end: number = Math.min(e.endTime - this.timeRange[0], this.timeRange[1])
      //   if (end >= start) {
      //     line = fill(line, 1, start, end + 1)
      //   }
      // }
      // this.timelines.push(line)

      let infectEventsLine: number[][] = []
      let start: number = 0
      let end: number = 0
      let mini: number = 0
      while (infectEvents[mini].startTime < timeRange[0]) {
        mini += 1
        if (mini >= infectEvents.length) {
          this.infectEventsLines.push([[0, timeRange[1] - timeRange[0] + 1, timeRange[0], timeRange[1], 0]])
          return
        }
      }
      let nextStart: number = Math.min(infectEvents[mini].startTime, timeRange[1] + 1)
      if (nextStart - 1 >= 0) {
        infectEventsLine.push([0, nextStart - timeRange[0], timeRange[0], nextStart - 1, 0])
      }
      for (let i = mini; i < eLen; i++) {
        const e = infectEvents[i]
        if (e.startTime > timeRange[1]) {
          break
        }
        if (e.endTime < timeRange[0]) {
          continue
        }
        start = Math.max(e.startTime, timeRange[0])
        end = Math.min(e.endTime, timeRange[1])
        if (end >= start) {
          infectEventsLine.push([1, end - start + 1, start, end, 0])
        }
        if (end >= timeRange[1]) {
          break
        }
        if (i + 1 < eLen) {
          const nextEvent: InfectEvent = infectEvents[i + 1]
          nextStart = Math.min(nextEvent.startTime, timeRange[1] + 1)
          if (end + 1 <= nextStart - 1) {
            infectEventsLine.push([0, nextStart - 1 - (end + 1) + 1, end + 1, nextStart - 1, 0])
          }
        }
      }
      this.infectEventsLines.push(infectEventsLine)
    })
    const influenceList: FullInfluenceInfo[] = this.timeWindowInfo.tree.FullInfluenceInfoList
    this.infectEventsLines.forEach((line: number[][], i: number) => {
      line.forEach((cell: number[], j: number) => {
        if (cell[0] === 1) {
          if (typeof find(influenceList, { LocationIDi: this.locationsId[i], StartTimeI: cell[2] }) === 'object' ||
          typeof find(influenceList, { LocationIDj: this.locationsId[i], StartTimeJ: cell[2] }) === 'object') {
            this.infectEventsLines[i][j][4] = 1
          } else {
            this.infectEventsLines[i][j][4] = 0.2
          }
        }
      })
    })
    const parent: {[key: string]: string} = {}
    const child: {[key: string]: string[]} = {}
    influenceList.forEach((influence: FullInfluenceInfo) => {
      const keyI: string = `${influence.LocationIDi},${influence.StartTimeI}`
      const keyJ: string = `${influence.LocationIDj},${influence.StartTimeJ}`
      parent[keyJ] = keyI
      if (keyI in child) {
        child[keyI].push(keyJ)
      } else {
        child[keyI] = []
      }
    })

    // draw a line to connect two cells when they has same parent
    for (let i = 0; i < this.infectEventsLines.length; i++) {
      const line: number[][] = this.infectEventsLines[i]
      for (let j = 0; j < line.length; j++) {
        if (line[j][0] === 0 || line[j][4] === 0.2) {
          continue
        }
        if (j + 2 >= line.length) {
          break
        }
        if (line[j + 2][4] !== 1) {
          continue
        }
        const keyI: string = `${this.locationsId[i]},${line[j][2]}`
        const keyJ: string = `${this.locationsId[i]},${line[j + 2][2]}`
        if (parent[keyI] === parent[keyJ]) {
          this.infectEventsLines[i][j + 1][0] = 2
          this.infectEventsLines[i][j + 1][4] = 1
        }
      }
    }
  }
}
</script>

<style scoped lang="scss">
.time-line-box-container {
  position: relative;
  display: inline-block;
  width: 200px;
  height: calc(100% - 18px);
  margin: 9px 0px 9px 9px;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
  transition: all 0.3s cubic-bezier(.25,.8,.25,1);

  &:hover {
    box-shadow: 0 10px 20px rgba(0,0,0,0.25), 0 6px 6px rgba(0,0,0,0.22);
  }

  .head {
    height: 50px;
    width: calc(100% - 5px);
    padding-left: 5px;
    overflow: hidden;
    text-overflow: ellipsis;
    .time-txt {
      width: 100%;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .timeline-box {
    padding: 0 8px;
    height: calc(100% - 50px);
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    align-items: center;
    position: relative;

    .line {
      padding-left: 5px;
      padding-right: 5px;
      width: calc(100% - 10px);
      height: 8px;
      display: flex;
      align-items: center;
      background-color: #f0f0f0;
      border-radius: 4px;

      .cell {
          height: 100%;
          flex: 1;
      }

      .cell-event {
        border-radius: 4px;
        background-color: #777;
      }

      /*.first-cell {*/
      /*  border-bottom-left-radius: 4px;*/
      /*  border-top-left-radius: 4px;*/
      /*}*/

      /*.last-cell {*/
      /*  border-bottom-right-radius: 4px;*/
      /*  border-top-right-radius: 4px;*/
      /*}*/
    }
  }
}
</style>
