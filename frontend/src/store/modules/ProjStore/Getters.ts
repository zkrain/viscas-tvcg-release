import { GetterTree } from 'vuex'

import { IRootState } from '../../RootState'
import { IProjStoreState } from './State'

import { filter, keyBy, forEach } from 'lodash'
import { hexToRgb } from '@/utils/draw'
import * as turf from '@turf/turf'

export const getters: GetterTree<IProjStoreState, IRootState> = {
  InfluencesWithIDListSlice: (state: IProjStoreState): InfluencesWithID[] => {
    return state.cascadingLinks.map((link: number[]) => {
      const locationIDi = link[0]
      const locationIDj = link[1]
      const edge = state.influencesDict[locationIDi][locationIDj]
      const edgeSample = filter(edge, (e: Influence) => e.StartTimeI < state.timeLength)
      return {
        influences: edgeSample,
        locationIDi,
        locationIDj
      }
    })
  },
  filteredIDDict: (state: IProjStoreState): {[key: number]: number} => {
    return keyBy(state.filteredIDs)
  },
  locationDict: (state: IProjStoreState): {[locationID: number]: ILocation} => {
    const locationDict: {[locationID: number]: ILocation} = {}
    state.allLocations.forEach((location: ILocation) => {
      locationDict[location.rid] = location
    })
    return locationDict
  },
  embedPointDict: (state: IProjStoreState): {[locationID: number]: ProjPoint} => {
    const embedPointDict: {[locationID: number]: ProjPoint} = {}
    state.embedPoints.forEach((embedPoint) => {
      embedPointDict[embedPoint.rid] = embedPoint
    })
    return embedPointDict
  },
  // 所有点的范围
  locationRange: (state: IProjStoreState): ILocationRange => {
    const tRange: ILocationRange = { minLat: 90, maxLat: -90, minLng: 180, maxLng: -180 }
    state.allLocations.forEach((location) => {
      tRange.minLat = Math.min(tRange.minLat, location.lat)
      tRange.maxLat = Math.max(tRange.maxLat, location.lat)
      tRange.minLng = Math.min(tRange.minLng, location.lng)
      tRange.maxLng = Math.max(tRange.maxLng, location.lng)
    })
    return tRange
  },

  // 饼状图中小扇形的缩放比例
  pieChartInnerRate: (state: IProjStoreState): number => {
    return state.pieChartInnerRate
  },

  edgeColorDict: (state: IProjStoreState): {[edgeKey: string]: [string, string]} => {
    const dict: {[edgeKey: string]: [string, string]} = {}
    state.cascadingLinks.forEach((link, i) => {
      const edgeKey: string = `${link[0]},${link[1]}`
      const rgb = hexToRgb(state.colors[i]) as {r: number, g: number, b: number}
      dict[edgeKey] = [`rgba(${rgb.r},${rgb.g},${rgb.b},0.15)`, state.colors[i]]
    })
    return dict
  },

  locationColorDict: (state: IProjStoreState): {[locationID: number]: string} => {
    const name = state.currentDataSetConfig.name
    const dict: {[locationID: number]: string} = {}
    let i = 0
    state.cascadingLinks.forEach((link) => {
      const locationIDi = link[0]
      const locationIDj = link[1]
      if (!(locationIDi in dict)) {
        dict[locationIDi] = state.locationColors[name][i]
        i++
      }
      if (!(locationIDj in dict)) {
        dict[locationIDj] = state.locationColors[name][i]
        i++
      }
    })
    return dict
  },

  // {[edgeKey: string (e.g., locationI,locationJ)]: normalized distance [0.2~1] }
  edgeLengthDict: (state: IProjStoreState): {[key: string]: number} => {
    const dict: {[key: string]: number} = {}
    if (state.cascadingLinks.length === 0) {
      return dict
    }
    let minDistance: number = 999999999
    let maxDistance: number = 0
    state.cascadingLinks.forEach((edge: number[]) => {
      const startLocation: ILocation = state.allLocationsInfoDict[edge[0]].location
      const endLocation: ILocation = state.allLocationsInfoDict[edge[1]].location
      const distance: number = turf.distance(turf.point([startLocation.lng, startLocation.lat]),
        turf.point([endLocation.lng, endLocation.lat]))
      minDistance = Math.min(minDistance, distance)
      maxDistance = Math.max(maxDistance, distance)
      dict[`${edge[0]},${edge[1]}`] = distance
    })
    if (minDistance === maxDistance) {
      dict[`${state.cascadingLinks[0][0]},${state.cascadingLinks[0][1]}`] = 1
      return dict
    }
    state.cascadingLinks.forEach((edge: number[]) => {
      const key: string = `${edge[0]},${edge[1]}`
      dict[key] = 0.2 + 0.8 * (dict[key] - minDistance) / (maxDistance - minDistance)
    })
    return dict
  },

  selectedTreesIdDictByATrees: (state: IProjStoreState): {[key: number]: boolean} => {
    const dict: {[key: number]: boolean} = {}
    forEach(state.trees, (t) => {
      dict[t.id] = false
    })
    forEach(state.selectedTreesId, (id) => {
      dict[id] = true
    })
    return dict
  },

  filteredTreesIndices: (state: IProjStoreState): {[index: number]: boolean} => {
    const indices: {[index: number]: boolean} = {}
    forEach(state.trees, (t) => {
      let flag = true
      forEach(state.selectedEdge, (value, edgeKey) => {
        if (!(edgeKey in t.involvedEdgeDict)) {
          flag = false
        }
      })
      if (flag) {
        indices[t.id] = true
      }
    })
    return indices
  },

  filteredTrees: (state: IProjStoreState): PropagationTree[] => {
    return filter(state.trees, (t) => {
      let flag = true
      forEach(state.selectedEdge, (value, edgeKey) => {
        if (!(edgeKey in t.involvedEdgeDict)) {
          flag = false
        }
      })
      return flag
    })
  },

  locationsScreenBound: (state: IProjStoreState): number[] => {
    const xRange: [number, number] = [999999, -999999]
    const yRange: [number, number] = [999999, -999999]
    for (let key in state.cascadingPointsScreenPosition) {
      const pt: [number, number] = state.cascadingPointsScreenPosition[key]
      xRange[0] = Math.min(xRange[0], pt[0])
      xRange[1] = Math.max(xRange[1], pt[0])
      yRange[0] = Math.min(yRange[0], pt[1])
      yRange[1] = Math.max(yRange[1], pt[1])
    }
    return [...xRange, ...yRange]
  }
}
