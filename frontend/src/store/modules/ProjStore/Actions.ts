import Vue from 'vue'
import { ActionContext, ActionTree, Action } from 'vuex'
import { MUTATIONS } from './MutationTypes'
import { IRootState } from '../../RootState'
import { IProjStoreState } from './State'
import EventBus from '../../../eventBus'
import { zipWith, join, reduce, sortBy, flatten, forEach, dropWhile, uniq } from 'lodash'
import * as d3 from 'd3'
import hull from 'hull.js'

export const actions: ActionTree<IProjStoreState, IRootState> = {
  async setDataSetConfig (ctx: ActionContext<IProjStoreState, IRootState>, dataSetConfig: IDataSetConfig): Promise<void> {
    ctx.commit(MUTATIONS.SET_DATA_SET_CONFIG, dataSetConfig)
    ctx.commit(MUTATIONS.SET_INFERENCE_RESULT, { influencesDict: [], cascadingLinks: [], trees: [] })
  },

  async fetchDataset (ctx: ActionContext<IProjStoreState, IRootState>, filename: string): Promise<void> {
    // embeddings
    // locations
    // events (not all will be fetched)
    console.log('changing Backend Dataset')
    await ctx.dispatch('changeBackendDataset', filename) // changeBackendDataset 在之后也是会用到的，所以抽出来

    // 后台换完数据了，是不是就不用传filename过去了？ to be checked
    // re-get newly updated data
    console.log('Backend Dataset has been changed')
    await ctx.dispatch('getDistribution')
    await ctx.dispatch('getAllLocations')
    await ctx.dispatch('getEmbeddings')
    await ctx.dispatch('getAllLocationsInfo', ctx.state.allIDs)
    await ctx.dispatch('getAllLocationsDataDict')
    // ctx.dispatch('getAllID')
  },

  async changeBackendDataset (ctx: ActionContext<IProjStoreState, IRootState>, filename: string): Promise<void> {
    const resp = await Vue.http.post(`${ctx.state.url}changeBackendDataset`, { filename })
    const code = resp.data.Code
    if (code === 200) {
      console.log('Error')
    } else {
      console.log('Done')
    }
  },

  async getEmbeddings (ctx: ActionContext<IProjStoreState, IRootState>): Promise<void> {
    const resp = await Vue.http.get(`${ctx.state.url}getEmbeddings`)
    const embeddings: ProjPoint[] = resp.data.Embeddings
    const allIDs: number[] = embeddings.map((embedding) => embedding.rid)
    ctx.commit(MUTATIONS.SET_EMBEDDINGS, embeddings)
    ctx.commit(MUTATIONS.SET_ALL_IDS, allIDs)
  },

  async getDistribution (ctx: ActionContext<IProjStoreState, IRootState>): Promise<void> {
    const distributionResp = await Vue.http.post(`${ctx.state.url}getDistribution`)
    const distribution = distributionResp.data.Distribution
    const maxCount = distributionResp.data.MaxCount
    const minCount = distributionResp.data.MinCount
    ctx.commit(MUTATIONS.SET_DISTRIBUTION, {
      distribution, maxCount, minCount
    })
  },

  // async getAllID (ctx: ActionContext<IProjStoreState, IRootState>): Promise<void> {
  //   const resp = await Vue.http.post(`${ctx.state.url}getAllID`)
  //   const allIDs = <number[]>resp.data.IDs

  //   ctx.commit(MUTATIONS.SET_ALL_IDS, allIDs)
  // },

  async filterIDsByCount (ctx: ActionContext<IProjStoreState, IRootState>, payload: {vLow: number, vHigh: number}): Promise<void> {
    const resp1 = await Vue.http.post(`${ctx.state.url}filterIDsByCount`, {
      vLow: Math.ceil(payload.vLow),
      vHigh: Math.ceil(payload.vHigh)
    })
    const filteredIDs = resp1.data.IDs
    ctx.commit(MUTATIONS.SET_FILTERED_IDS, filteredIDs)
  },

  async getAllLocations (ctx: ActionContext<IProjStoreState, IRootState>): Promise<void> {
    const locationsResp = await Vue.http.post(`${ctx.state.url}getAllLocations`)
    const locations: ILocation[] = locationsResp.data.Locations
    ctx.commit(MUTATIONS.SET_ALL_LOCATIONS, locations)
  },

  async getAllLocationsInfo (ctx: ActionContext<IProjStoreState, IRootState>, locationIDs: number): Promise<void> {
    const locationsResp = await Vue.http.post(`${ctx.state.url}getLocations`, { locationIDs })
    const locations: ILocation[] = locationsResp.data.Locations
    const locationsStyle: ILocationStyle[] = locationsResp.data.LocationsStyle
    const locationsEventCount: number[] = locationsResp.data.LocationsEventCount
    const locationInfoDict: {[key: number]: ILocationInfo} = {}
    zipWith(locations, locationsEventCount, locationsStyle, (location, totalEventCount, locationStyle) => {
      const locationInfo: ILocationInfo = { location, totalEventCount, locationStyle }
      locationInfoDict[location.rid] = locationInfo
    })
    ctx.commit(MUTATIONS.SET_LOCATIONS_INFO_DICT, locationInfoDict)
  },

  async getAllLocationsDataDict (ctx: ActionContext<IProjStoreState, IRootState>) {
    const resp = await Vue.http.post(`${ctx.state.url}getLocationAverage`)
    const locationsAverageDataDict: {[key: number]: [number, number][]} = resp.body.averageEventDict
    ctx.commit(MUTATIONS.SET_ALL_LOCATIONS_AVERAGE_DATA_DICT, locationsAverageDataDict)
  },

  async getCorrelatedLocations (ctx: ActionContext<IProjStoreState, IRootState>,
    params: { clusterId: number, locationIDs: number[]}): Promise<void> {
    console.log('getCorrelatedLocations')
    const locationIDs = params.locationIDs
    const locationsResp = await Vue.http.post(`${ctx.state.url}getLocations`, { locationIDs })
    const locations: ILocation[] = locationsResp.data.Locations
    const locationsStyle: ILocationStyle[] = locationsResp.data.LocationsStyle
    const locationsEventCount: number[] = locationsResp.data.LocationsEventCount
    const locationsInfo: ILocationInfo[] = zipWith(locations, locationsEventCount, locationsStyle, (location, totalEventCount, locationStyle) => {
      return { location, totalEventCount, locationStyle }
    })
    ctx.commit(MUTATIONS.SET_LOCATIONS_INFO, locationsInfo)
  },

  async getInfectEventsDict (ctx: ActionContext<IProjStoreState, IRootState>, locationIDs: number[]) {
    const locationsResp = await Vue.http.post(`${ctx.state.url}getLocationInfectEvents`, { locationIDs })
    ctx.commit(MUTATIONS.SET_CLUSTER_INFECT_EVENTS_DICT, locationsResp.body.infectEventsDict)
  },

  async getEventsOfCorrelatedLocations (ctx: ActionContext<IProjStoreState, IRootState>, locationIDs: number[]): Promise<void> {
    const eventsResp = await Vue.http.post(`${ctx.state.url}getEventsOfCorrelatedLocations`, { locationIDs })
    const eventsList: number[][] = eventsResp.data.Events

    const eventsDict: {[locationID: number]: number[]} = {}
    const eventsWithIDList: IEventsWithID[] = []
    locationIDs.forEach((locationID: number, index: number) => {
      eventsDict[locationID] = eventsList[index]
      eventsWithIDList.push({
        locationID: locationID,
        events: eventsList[index]
      })
    })

    ctx.commit(MUTATIONS.SET_EVENTS_DICT, eventsDict)
    ctx.commit(MUTATIONS.SET_EVENTS_WITH_ID_LIST, eventsWithIDList)
  },

  async inferCascadingPattern (ctx: ActionContext<IProjStoreState, IRootState>, { locationIDs, k, tw }: { locationIDs: number[], k: number, tw: number }): Promise<void> {
    await ctx.dispatch('clearSelectedEdge')
    await ctx.dispatch('clearSelectedAtree')

    const resp = await Vue.http.post(`${ctx.state.url}inferCascadingPattern`, { locationIDs, k, tw })
    const cascadingLinks: [number, number, number][] = resp.data.links.map((link: [number, number], i: number) => [link[0], link[1], i])
    const influencesDict: InfluencesDict = resp.data.influencesDict
    const trees: PropagationTree[] = resp.data.trees
    // const influencesDictContraryInG: {[linkKey: string]: Influence[]} = resp.data.influencesDictContraryInG
    // const influencesDictContraryNotInG: {[linkKey: string]: Influence[]} = resp.data.influencesDictContraryNotInG

    const involvedLocations: number[] = []
    cascadingLinks.forEach((link) => {
      const locationIDi = link[0]
      const locationIDj = link[1]
      if (involvedLocations.indexOf(locationIDi) < 0) {
        involvedLocations.push(locationIDi)
      }
      if (involvedLocations.indexOf(locationIDj) < 0) {
        involvedLocations.push(locationIDj)
      }
    })

    // aggregate the p by edge
    const edgeProbabilityDistribution: {[edgeKey: string]: number[]} = {}

    const totalInfluenceDuration = ctx.state.currentDataSetConfig.totalInfluenceDuration
    const timeWindows: ITimeWindowInfo[] = []

    // for each tree, obtain their information
    trees.forEach((tree: PropagationTree, index: number) => {
      const durationDict: {[locationID: number]: {
        [startTime: number]: number
      }} = {}
      const involvedEdgeDict: {[edgeKey: string]: SimplifiedInfluenceInfo} = {}
      const involvedLocationDict: {[locationID: number]: number} = {}
      tree.FullInfluenceInfoList.forEach((fullInfluenceInfo) => {
        const locationIDi = fullInfluenceInfo.LocationIDi
        const locationIDj = fullInfluenceInfo.LocationIDj
        const startTimeI = fullInfluenceInfo.StartTimeI
        const startTimeJ = fullInfluenceInfo.StartTimeJ
        const durationI = fullInfluenceInfo.DurationI + 1
        const durationJ = fullInfluenceInfo.DurationJ + 1
        const p = fullInfluenceInfo.P
        const edgeKey = locationIDi + ',' + locationIDj

        // count p in edgeProbabilityDistribution for each edge
        if (!(edgeKey in edgeProbabilityDistribution)) {
          const probabilityDistribution = new Array(10).fill(0)
          edgeProbabilityDistribution[edgeKey] = probabilityDistribution
        }
        const pLowerBound = ctx.state.currentDataSetConfig.pLowerBound
        const pIndex = Math.floor((p - pLowerBound) / ((1 - pLowerBound) / 10))
        edgeProbabilityDistribution[edgeKey][pIndex]++

        if (!(locationIDi in durationDict)) {
          durationDict[locationIDi] = {}
        }
        if (!((+locationIDj) in durationDict)) {
          durationDict[(+locationIDj)] = {}
        }
        durationDict[(+locationIDj)][startTimeJ] = durationJ
        durationDict[locationIDi][startTimeI] = durationI

        // organize potential causes 1
        let pso = fullInfluenceInfo.OtherPotentialCauses.map(c => [c.P, c.Duration] as [number, number]).sort((a, b) => a[0] - b[0])
        pso = pso.map((po) => {
          if (po[0] > p) {
            return [(po[0] - p) / (1 - p), po[1]]
          } else {
            return [-(p - po[0]) / (p - 0.3), po[1]]
          }
        })

        // organize potential cause
        const psoDict: {[locationID: number]: PotentialCause} = {}
        fullInfluenceInfo.OtherPotentialCauses.forEach(c => {
          psoDict[c.LocationID] = c
          // if ((c.Duration + 1 >= durationI) && (c.P > p)) {
          //   console.log(locationIDi, locationIDj, 'other: ', c.Duration + 1 >= durationI, c.Duration + 1, c.P, 'self', durationI, p)
          // }
        })
        const potentialCausesByLocation: PotentialCause[] = []
        locationIDs.forEach((lid) => {
          if (lid !== locationIDj) {
            if (lid in psoDict) {
              let relativeP = psoDict[lid].P
              if (relativeP > p) {
                relativeP = (relativeP - p) / (1 - p)
              } else {
                relativeP = -(p - relativeP) / (p - 0.3)
              }
              potentialCausesByLocation.push({
                LocationID: lid,
                StartTime: psoDict[lid].StartTime,
                P: relativeP,
                Duration: psoDict[lid].Duration
              })
            } else if (lid !== locationIDi) {
              potentialCausesByLocation.push({
                LocationID: lid,
                StartTime: -1,
                P: 0,
                Duration: 0
              })
            }
          }
        })

        // organize potential causes 4 // 暂时只除去自己和自己的父节点，pattern里的其他节点保留
        let rankOfDuration = locationIDs.length - 2
        let finalP = 0
        let nLagre = 0
        locationIDs.forEach((lid) => {
          if (lid !== locationIDj) {
            if (lid in psoDict) {
              if (p < psoDict[lid].P) {
                finalP += (p - psoDict[lid].P)
                if (durationI > psoDict[lid].Duration) {
                  rankOfDuration--
                }
                nLagre++
              }
            }
          }
        })
        rankOfDuration /= (locationIDs.length - 2)
        finalP = nLagre === 0 ? p : (finalP / nLagre / p + p)
        // console.log(locationIDs.length - 2, rankOfDuration, finalP)

        const sti = startTimeI
        let reti = durationI
        let rstj = startTimeJ - startTimeI
        let retj = startTimeJ - startTimeI + durationJ
        reti = reti >= totalInfluenceDuration ? (totalInfluenceDuration - 1) : reti
        rstj = rstj >= totalInfluenceDuration ? (totalInfluenceDuration - 1) : rstj
        retj = retj >= totalInfluenceDuration ? (totalInfluenceDuration - 1) : retj

        // get durationLinkParams
        let durationLinks: [number, number, number][] = []
        if (retj <= reti) {
          durationLinks = [[0, rstj, 0], [rstj, retj, -1], [retj, reti, 0]]
        } else if (reti < rstj) {
          durationLinks = [[0, reti, 0], [rstj, retj, 1]]
        } else {
          durationLinks = [[0, rstj, 0], [rstj, reti, -1], [reti, retj, 1]]
        }

        // get durationLinkParams 2
        let durationLinksAlignByRstj: [number, number, number][] = []
        if (retj <= reti) {
          durationLinksAlignByRstj = [[tw - rstj, tw, 0], [tw, tw + retj - rstj, -1], [tw + retj - rstj, tw + reti - rstj, 0]]
        } else if (reti < rstj) {
          durationLinksAlignByRstj = [[tw - rstj, tw - rstj + reti, 0], [tw, tw + retj - rstj, 1]]
        } else {
          durationLinksAlignByRstj = [[tw - rstj, tw, 0], [tw, tw + reti - rstj, -1], [tw + reti - rstj, tw + retj - rstj, 1]]
        }

        if (!(edgeKey in involvedEdgeDict) || p > involvedEdgeDict[edgeKey].p) {
          // 如果边不在字典中
          // 或者p的值大于里面的值
          involvedEdgeDict[edgeKey] = {
            sti,
            reti,
            rstj,
            retj,
            p,
            pso,
            duration: durationI,
            csoOrderedByLocation: potentialCausesByLocation,
            opcs: fullInfluenceInfo.OtherPotentialCauses,
            durationLinks,
            durationLinksAlignByRstj,
            rankOfDuration,
            finalP
          }
        }
      })

      let n = 0
      let averageDurationInTree = 0
      forEach(involvedLocations, (locationID) => {
        let averageDurationInLocation = 0
        let nn = 0
        forEach(durationDict[locationID], (duration, startTime) => {
          n++
          nn++
          averageDurationInTree += duration
          averageDurationInLocation += duration
        })
        involvedLocationDict[+locationID] = nn > 0 ? averageDurationInLocation / nn : 0
      })
      averageDurationInTree /= n

      const durations: [number, number][] = []
      forEach(involvedLocationDict, (duration, locationID) => {
        durations.push([duration, +locationID])
      })

      tree.averageDuration = averageDurationInTree
      tree.durationDict = durationDict
      tree.durations = durations
      tree.involvedEdges = Object.keys(involvedEdgeDict)
        .map(edgeKey => edgeKey.split(',')
          .map(locationID => (+locationID)) as [number, number])
      tree.involvedEdgeDict = involvedEdgeDict
      tree.id = index
      timeWindows.push({ timeRange: [tree.TreeStartTime, tree.TreeEndTime], tree })
    })

    // normalize edgeProbabilityDistribution
    forEach(edgeProbabilityDistribution, (probabilityDistribution, edgeKey) => {
      const maxNumber = Math.max(...probabilityDistribution)
      edgeProbabilityDistribution[edgeKey] = probabilityDistribution.map(number => number / maxNumber)
    })

    // Classify the trees by their structure
    const mapEdgeKey2Index: {[edgeKey: string]: number} = {}
    cascadingLinks.forEach((link, i) => {
      mapEdgeKey2Index[`${link[0]},${link[1]}`] = i
    })
    const treeCollection: {[topologyKey: string]: PropagationTree[]} = {}
    trees.forEach((tree, treeIndex) => {
      const edgeKeys = Object.keys(tree.involvedEdgeDict)
      const indices = edgeKeys.map((edgeKey) => {
        const index = mapEdgeKey2Index[edgeKey]
        return index
      })
      indices.sort()
      const topologyKey = join(indices, '_')
      if (topologyKey in treeCollection) {
        treeCollection[topologyKey].push(tree)
      } else {
        treeCollection[topologyKey] = [tree]
      }
    })

    // for each sets of trees that have the same structure, we calculate their information and store them in aggregatePropagationTrees
    const aggregatePropagationTrees: AggregatePropagationTree[] = []
    let aid: number = 0
    forEach(treeCollection, (treesWithSameTopo, topologyKey) => {
      const involvedEdgeDictInAggregation: {[edgeKey: string]: SimplifiedInfluenceInfo[]} = {}

      // all edgeKeys, like "locationIDi + ',' + locationIDj", are stored in the keys of treesWithSameTopo[0].involvedEdgeDict
      // initialize involvedEdgeDictInAggregation
      forEach(treesWithSameTopo[0].involvedEdgeDict, (edgeInfo, edgeKey) => { // edgeInfo : number[5]
        // for each edge
        involvedEdgeDictInAggregation[edgeKey] = []
      })

      treesWithSameTopo.forEach((tree) => {
        forEach(tree.involvedEdgeDict, (edgeInfo, edgeKey) => {
          involvedEdgeDictInAggregation[edgeKey].push(edgeInfo)
        })
      })

      for (let i = 0; i < treesWithSameTopo.length; i++) {
        treesWithSameTopo[i].parentAid = aid
      }

      forEach(involvedEdgeDictInAggregation, async (value, edgeKey) => {
        // get influence vertices
        // const vertices = involvedEdgeDictInAggregation[edgeKey]
        //   .map((influence) => [influence.reti, influence.rstj, influence.retj] as [number, number, number])
        // get influence vertices 2
        // const vertices = involvedEdgeDictInAggregation[edgeKey]
        //   .map((influence) => [Math.round(influence.finalP * 8), -influence.rstj * 2, influence.reti - influence.rstj, influence.retj - influence.rstj] as [number, number, number, number])
        // get influence vertices 3
        const vertices = involvedEdgeDictInAggregation[edgeKey]
          .map((influence) => [-influence.rstj * 1, 1 * (influence.reti - influence.rstj), (influence.retj - influence.rstj)] as [number, number, number])
        const respHamiltonianInfluence = await Vue.http.post(`${ctx.state.url}hamiltonianWalkOptimizerInfluence`, { vertices })
        const influenceInfoOrder: number[] = respHamiltonianInfluence.data.order
        const influenceInfoListOrdered: SimplifiedInfluenceInfo[] = []
        influenceInfoOrder.forEach(i => {
          influenceInfoListOrdered.push(involvedEdgeDictInAggregation[edgeKey][i])
        })
        // const influenceInfoListOrdered = involvedEdgeDictInAggregation[edgeKey].influenceInfoList.sort((inf1, inf2) => (inf1.rstj - inf2.rstj))
        involvedEdgeDictInAggregation[edgeKey] = influenceInfoListOrdered
      })

      aggregatePropagationTrees.push({
        trees: treesWithSameTopo,
        id: aid,
        selected: false,
        // height: 16 * Math.log(treesWithSameTopo.length) + 8,
        height: treesWithSameTopo.length * 4 + 6,
        involvedEdges: treesWithSameTopo[0].involvedEdges,
        involvedEdgeDictInAggregation
      })
      aid += 1
    })

    // generate vertices of hamiltonian graph by cascading link
    const verticesEdge: number[][] = []
    cascadingLinks.forEach((link) => {
      const edgeKey = `${link[0]},${link[1]}`
      const vertex: number[] = []
      aggregatePropagationTrees.forEach((at, i) => {
        if (edgeKey in at.involvedEdgeDictInAggregation) {
          vertex.push(Math.round(at.height))
        } else {
          vertex.push(0)
        }
      })
      verticesEdge.push(vertex)
    })
    const respHamiltonianOptimal = await Vue.http.post(`${ctx.state.url}hamiltonianWalkOptimizerEdge`, { vertices: verticesEdge })
    const orderEdge: number[] = respHamiltonianOptimal.data.order
    const cascadingLinksOrdered: [number, number, number][] = []
    orderEdge.forEach(i => {
      cascadingLinksOrdered.push(cascadingLinks[i])
    })

    // finally, set in state
    let locationsId: number[] = []
    cascadingLinks.forEach((t) => {
      locationsId.push(t[0], t[1])
    })
    locationsId = uniq(locationsId)
    await ctx.dispatch('getInfectEventsDict', locationsId)
    ctx.commit(MUTATIONS.SET_INFERENCE_RESULT, {
      influencesDict,
      cascadingLinks: cascadingLinksOrdered,
      trees,
      influencesContrary: {},
      aggregatePropagationTrees,
      edgeProbabilityDistribution,
      timeWindows
    })
    console.log('aggregatePropagationTrees', aggregatePropagationTrees)
    console.log(cascadingLinks)
  },

  async clearCascadingPattern (ctx: ActionContext<IProjStoreState, IRootState>): Promise<void> {
    ctx.commit(MUTATIONS.SET_INFERENCE_RESULT, { influencesDict: [], cascadingLinks: [], trees: [] })
  },

  async reduceNoiseForEvents (ctx: ActionContext<IProjStoreState, IRootState>, payload: { window: number, threshold: number }): Promise<void> {
    const resp = await Vue.http.post(`${ctx.state.url}reduceNoiseForEvents`, payload)
    const eventsList = resp.data.Events
    console.log(eventsList)

    ctx.commit(MUTATIONS.UPDATE_EVENTS_LIST, eventsList)
    console.log('events filtered/updated')
  },

  // async setClusterResult (ctx: ActionContext<IProjStoreState, IRootState>, payload: {
  //   colorMapForEachLocation: {[rid: number]: string},
  //   clustersWithContext: ICluster[],
  //   maxValues: { maxStartTimeNumber: number, maxRawEventNumber: number, maxLengthNumber: number }
  // }): Promise<void> {
  //   ctx.commit(MUTATIONS.SET_CLUSTER_RESULT, payload)
  // },

  async setProjectionCanvasSize (ctx: ActionContext<IProjStoreState, IRootState>, payload: {width: number, height: number}) {
    ctx.commit(MUTATIONS.SET_PROJECTION_CANVAS_SIZE, payload)
  },

  async stclustering (ctx: ActionContext<IProjStoreState, IRootState>, payload: {
    EpsT: number,
    EpsS: number,
    MinPts: number,
    deltaE: number
  }): Promise<void> {
    const locationIDs: number[] = ctx.state.filteredIDs
    const embedPointDict = ctx.getters.embedPointDict
    const locationDict: {[locationID: number]: ILocation} = ctx.getters.locationDict

    const resp = await Vue.http.post(`${ctx.state.url}stclustering`, { ...payload, locationIDs })
    const averageEventCount: number[] = resp.data.AverageEventCount
    const clusterEventCount: number[][] = resp.data.ClusterEventCount
    const averageInfectEventCount = resp.data.AverageInfectEventCount
    const clusters = resp.data.Clusters
    const noise = resp.data.Noise
    const clusterInformationList: ClusterInformation[] = resp.data.ClusterInformationList
    const spatialDistributionOfClusters: SpatialDistributionOfClusters = resp.data.SpatialDistributionOfClusters
    const distributions: number[][] = spatialDistributionOfClusters.DistributionList

    console.log(distributions)

    const colorMapForEachLocation: {[rid: number]: string} = {}
    const colors = d3.scaleSequential(d3.interpolateRainbow)
      .domain([0, clusters.length])

    const clustersWithContext: ICluster[] = []
    let maxStartTimeNumber = 0
    let maxRawEventNumber = 0
    let maxLengthNumber = 0

    // calculate range of centers
    const clustersCenter: [number, number][] = []
    clusters.forEach((cRids: number[], ci: number) => {
      const center: [number, number] = [0, 0] // [lng, lat]
      cRids.forEach((rid: number) => {
        center[0] += locationDict[rid].lng
        center[1] += locationDict[rid].lat
      })
      center[0] /= cRids.length
      center[1] /= cRids.length
      clustersCenter.push(center)
    })

    // asign color to each cluster
    clusters.forEach((cRid: number[], ci: number) => {
      const clusterInformation = clusterInformationList[ci]
      const spatialDistribution = distributions[ci]

      // getMaxValue across cluster information
      clusterInformation.StartTimeNumberMatrix.forEach((vector) => {
        vector.forEach((v) => {
          if (v > maxStartTimeNumber) {
            maxStartTimeNumber = v
          }
        })
      })
      clusterInformation.RawEventNumberMatrix.forEach((vector) => {
        vector.forEach((v) => {
          if (v > maxRawEventNumber) {
            maxRawEventNumber = v
          }
        })
      })
      clusterInformation.LengthMatrix.forEach((vector) => {
        vector.forEach((v) => {
          if (v > maxLengthNumber) {
            maxLengthNumber = v
          }
        })
      })

      const cWithContext: ICluster = {
        color: colors(ci),
        rIDs: cRid,
        hull: '',
        hullPoints: [],
        clusterEventCount: [],
        clusterAverageEventCount: -1,
        clusterAverageInfectEventCount: -1,
        clusterInformation,
        spatialDistribution,
        spatialRelativeCenter: { rx: 0, ry: 0 }
      }

      // compute relative position of center in cluster-circle
      const geoRange = ctx.getters.locationRange
      const geoCenter: [number, number] = clustersCenter[ci]

      // 矩形内的坐标转为圆形(r=1)内的坐标
      const normedGeoCenter: [number, number] = [0, 0] // [lng, lat]
      normedGeoCenter[0] = (geoCenter[0] - geoRange.minLng) / (geoRange.maxLng - geoRange.minLng) - 0.5
      normedGeoCenter[1] = (geoCenter[1] - geoRange.minLat) / (geoRange.maxLat - geoRange.minLat) - 0.5
      const rRate: number = 2 * Math.sqrt(normedGeoCenter[0] * normedGeoCenter[0] + normedGeoCenter[1] * normedGeoCenter[1])
      cWithContext.spatialRelativeCenter = { rx: rRate * normedGeoCenter[0], ry: rRate * normedGeoCenter[1] }

      // compute hull and colorMapForEachLocation
      const XYs: number[][] = []

      cRid.forEach((rid: number) => {
        colorMapForEachLocation[rid] = colors(ci)
        XYs.push([embedPointDict[rid].left * ctx.state.projectionWidth, embedPointDict[rid].top * ctx.state.projectionHeight])
      })
      const hullPoints: [number, number][] = hull(XYs, 30)
      cWithContext.hullPoints = hullPoints
      // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
      cWithContext.hull = <string>d3.line()
        .x((d: [number, number]) => d[0])
        .y((d: [number, number]) => d[1])(hullPoints)
      clustersWithContext.push(cWithContext)
    })
    for (let i = 0; i < clustersWithContext.length; i++) {
      clustersWithContext[i].clusterAverageEventCount = averageEventCount[i]
      clustersWithContext[i].clusterAverageInfectEventCount = averageInfectEventCount[i]
      clustersWithContext[i].clusterEventCount = clusterEventCount[i]
    }

    // asign color to outliers
    noise.forEach((rid: number) => {
      colorMapForEachLocation[rid] = '#ccc'
    })

    ctx.commit(MUTATIONS.SET_CLUSTER_RESULT, {
      colorMapForEachLocation,
      clustersWithContext,
      maxValues: { maxStartTimeNumber, maxRawEventNumber, maxLengthNumber }
    })

    // ctx.dispatch('setClusterResult', {
    //   colorMapForEachLocation,
    //   clustersWithContext,
    //   maxValues: { maxStartTimeNumber, maxRawEventNumber, maxLengthNumber }
    // })
  },
  async setPieChartInnerRate (ctx: ActionContext<IProjStoreState, IRootState>, rate: number) {
    ctx.commit(MUTATIONS.SET_PIE_CHART_INNER_RATE, rate)
  },

  async setCascadePointsScreenPosition (ctx: ActionContext<IProjStoreState, IRootState>,
    positions: {[locationID: number]: [number, number]}) {
    ctx.commit(MUTATIONS.SET_CASCADING_POINTS_SCREEN_POSITION, positions)
  },

  async selectEdge (ctx: ActionContext<IProjStoreState, IRootState>, edgeKey: string) {
    ctx.commit(MUTATIONS.SET_SELECTED_EDGE, edgeKey)
  },

  async openEdge (ctx: ActionContext<IProjStoreState, IRootState>, edgeKey: string) {
    ctx.commit(MUTATIONS.SET_OPENED_EDGE, edgeKey)
  },

  async clearSelectedEdge (ctx: ActionContext<IProjStoreState, IRootState>) {
    ctx.commit(MUTATIONS.CLEAR_SELECTED_EDGE)
  },

  async selectTrees (ctx: ActionContext<IProjStoreState, IRootState>, trees: PropagationTree[]) {
    ctx.commit(MUTATIONS.XOR_SELECTED_TREES_ID, trees)
  },

  async selectAtree (ctx: ActionContext<IProjStoreState, IRootState>, ind: number) {
    ctx.commit(MUTATIONS.SELECT_ATREE, ind)
  },

  async clearSelectedAtree (ctx: ActionContext<IProjStoreState, IRootState>) {
    ctx.commit(MUTATIONS.CLEAR_SELECTED_ATREE)
  },

  async dispatchTimeWindows (ctx: ActionContext<IProjStoreState, IRootState>, timeWindows: ITimeWindowInfo[]) {
    ctx.commit(MUTATIONS.SET_TIME_WINDOWS, timeWindows)
  }
}
