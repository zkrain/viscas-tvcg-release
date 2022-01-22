import { MutationTree } from 'vuex'
import { findIndex, reject, without, forEach } from 'lodash'
import { MUTATIONS } from './MutationTypes'
import { IProjStoreState, state } from './State'
import Vue from 'vue'
import { Util } from 'leaflet'
import indexOf = Util.indexOf

// tslint:disable:function-name
export const mutations: MutationTree<IProjStoreState> = {
  [MUTATIONS.SET_DATA_SET_CONFIG] (state: IProjStoreState, dataSetConfig: IDataSetConfig) {
    state.currentDataSetConfig = dataSetConfig
    state.timeWindows = []
  },
  [MUTATIONS.SET_EMBEDDINGS] (state: IProjStoreState, points: ProjPoint[]) {
    state.embedPoints = points
  },
  [MUTATIONS.SET_INFERENCE_RESULT] (state: IProjStoreState,
    { influencesDict, cascadingLinks, trees,
      influencesContrary, aggregatePropagationTrees,
      edgeProbabilityDistribution, timeWindows }:
                                      { influencesDict: InfluencesDict, cascadingLinks: number[][],
                                        trees: PropagationTree[], influencesContrary: InfluencesDict,
                                        aggregatePropagationTrees: AggregatePropagationTree[],
                                        edgeProbabilityDistribution: {[edgeKey: string]: number[]},
                                        timeWindows: ITimeWindowInfo[]}) {
    state.cascadingLinks = cascadingLinks
    console.log('cascading links')
    console.log(state.cascadingLinks)
    state.timeWindows = []
    state.influencesDict = influencesDict // 可能会卡 再说 todo
    state.trees = trees
    state.influencesContrary = influencesContrary
    state.aggregatePropagationTrees = aggregatePropagationTrees
    state.edgeProbabilityDistribution = edgeProbabilityDistribution
    state.timeWindows = timeWindows
  },
  [MUTATIONS.UPDATE_EVENTS_LIST] (state: IProjStoreState, eventsList: number[][]) {
    state.eventsWithIDList.forEach((eventsWithID, i) => {
      eventsWithID.events = eventsList[i]
    })
  },
  [MUTATIONS.SET_ALL_IDS] (state: IProjStoreState, allIDs: number[]) {
    state.allIDs = allIDs

    // pre-set filter roads as all roads
    state.filteredIDs = allIDs
  },
  [MUTATIONS.SET_DISTRIBUTION] (state: IProjStoreState, distributionInfo: {
    distribution: number[]
    minCount: number
    maxCount: number
  }) {
    state.distribution = distributionInfo.distribution
    state.minCount = distributionInfo.minCount
    state.maxCount = distributionInfo.maxCount
  },
  [MUTATIONS.SET_FILTERED_IDS] (state: IProjStoreState, IDs: number[]) {
    state.filteredIDs = IDs
  },
  [MUTATIONS.SET_CLUSTER_RESULT] (state: IProjStoreState, payload: {
    colorMapForEachLocation: {[rid: number]: string},
    clustersWithContext: ICluster[],
    maxValues: IMaxValueAcrossClusters
  }) {
    state.colorMapForEachLocation = payload.colorMapForEachLocation
    state.clustersWithContext = payload.clustersWithContext
    state.maxValueAcrossClusters = payload.maxValues
  },
  [MUTATIONS.SET_PROJECTION_CANVAS_SIZE] (state: IProjStoreState, payload: {width: number, height: number}) {
    state.projectionWidth = payload.width
    state.projectionHeight = payload.height
  },
  [MUTATIONS.SET_EVENTS_DICT] (state: IProjStoreState, eventsDict: {[locationID: number]: number[]}) {
    state.eventsDict = eventsDict
  },
  [MUTATIONS.SET_EVENTS_WITH_ID_LIST] (state: IProjStoreState, eventsWithIDList: IEventsWithID[]) {
    state.eventsWithIDList = eventsWithIDList
  },
  [MUTATIONS.SET_ALL_LOCATIONS] (state: IProjStoreState, locations: ILocation[]) {
    state.allLocations = locations
  },
  [MUTATIONS.SET_PIE_CHART_INNER_RATE] (state: IProjStoreState, rate: number) {
    state.pieChartInnerRate = rate
  },
  [MUTATIONS.SET_CASCADING_POINTS_SCREEN_POSITION] (state: IProjStoreState,
    positions: {[locationID: number]: [number, number]}) {
    state.cascadingPointsScreenPosition = positions
  },
  [MUTATIONS.ALL_LOCATIONS_INFO_DICT_ADD] (state: IProjStoreState, params: {
    locationId: number, locationInfo: ILocationInfo
  }) {
    state.allLocationsInfoDict[params.locationId] = params.locationInfo
  },
  [MUTATIONS.SET_LOCATIONS_INFO_DICT] (state: IProjStoreState, locationInfoDict: {[key: number]: ILocationInfo}) {
    state.allLocationsInfoDict = locationInfoDict
  },
  [MUTATIONS.SET_LOCATIONS_INFO] (state: IProjStoreState, locationsInfo: ILocationInfo[]) {
    state.locationsInfo = locationsInfo
  },
  [MUTATIONS.SET_ALL_LOCATIONS_AVERAGE_DATA_DICT] (state: IProjStoreState, averageDataDict: {[key: number]: [number, number][]}) {
    state.allLocationsAverageDataDict = averageDataDict
  },

  [MUTATIONS.SET_OPENED_EDGE] (state: IProjStoreState, edgeKey: string) {
    if (edgeKey in state.openedEdge) {
      Vue.delete(state.openedEdge, edgeKey)
    } else {
      Vue.set(state.openedEdge, edgeKey, true)
    }
    console.log(state.openedEdge, 'openedEdge')
  },

  [MUTATIONS.SET_SELECTED_EDGE] (state: IProjStoreState, edgeKey: string) {
    if (edgeKey in state.selectedEdge) {
      Vue.delete(state.selectedEdge, edgeKey)
    } else {
      Vue.set(state.selectedEdge, edgeKey, true)
    }
    // function cmp (t1: AggregatePropagationTree, t2: AggregatePropagationTree): number {
    //   let okEdgeT1: number = 0
    //   let okEdgeT2: number = 0
    //   forEach(state.selectedEdge, (t:boolean, key: string) => {
    //     if (key in t1.involvedEdgeDictInAggregation && t) {
    //       okEdgeT1++
    //     }
    //     if (key in t2.involvedEdgeDictInAggregation && t) {
    //       okEdgeT2++
    //     }
    //   })
    //   return okEdgeT2 - okEdgeT1
    // }
    // state.aggregatePropagationTrees.sort(cmp)
  },
  [MUTATIONS.CLEAR_SELECTED_EDGE] (state: IProjStoreState) {
    const keys: string[] = Object.keys(state.selectedEdge)
    keys.forEach((key: string) => {
      Vue.delete(state.selectedEdge, key)
    })
  },

  [MUTATIONS.SET_CLUSTER_INFECT_EVENTS_DICT] (state: IProjStoreState, infectEventDict: {[key: number]: InfectEvent[]}) {
    state.clusterInfectEventsDict = infectEventDict
  },

  [MUTATIONS.SET_TIME_WINDOWS] (state: IProjStoreState, timeWindows: ITimeWindowInfo[]) {
    state.timeWindows = []
    timeWindows.forEach((t) => { // 为了触发v-for事件
      state.timeWindows.push(t)
    })
  },

  [MUTATIONS.XOR_SELECTED_TREES_ID] (state: IProjStoreState, trees: PropagationTree[]) {
    trees.forEach((t) => {
      if (indexOf(state.selectedTreesId, t.id) > -1) {
        state.selectedTreesId = without(state.selectedTreesId, t.id)
      } else {
        state.selectedTreesId.push(t.id)
      }
    })
  },

  [MUTATIONS.SELECT_ATREE] (state: IProjStoreState, aid: number) {
    const index: number = findIndex(state.aggregatePropagationTrees, { id: aid })
    if (index > -1) {
      const atree: AggregatePropagationTree = state.aggregatePropagationTrees[index]
      atree.selected = !atree.selected
      Vue.set(state.aggregatePropagationTrees, index, atree)
    }

    // push or pop(if existed) id of the aggregate propagation tree
    const ti: number = indexOf(state.selectedAggregatePropagationTreesId, aid)
    if (ti > -1) {
      state.selectedAggregatePropagationTreesId = without(state.selectedAggregatePropagationTreesId, aid)
    } else {
      state.selectedAggregatePropagationTreesId.push(aid)
    }
  },

  [MUTATIONS.CLEAR_SELECTED_ATREE] (state: IProjStoreState) {
    const ids: number[] = state.selectedAggregatePropagationTreesId
    ids.forEach((id: number) => {
      const index: number = findIndex(state.aggregatePropagationTrees, { id: id })
      if (index > -1) {
        const atree: AggregatePropagationTree = state.aggregatePropagationTrees[index]
        atree.selected = !atree.selected
        Vue.set(state.aggregatePropagationTrees, index, atree)
      }
    })
    state.selectedAggregatePropagationTreesId = []
    state.selectedTreesId = []
  }
}
