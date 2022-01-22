export interface IProjStoreState {
  currentDataSetConfig: IDataSetConfig
  url: string,
  embedPoints: ProjPoint[]
  cascadingLinks: number[][]
  edgeProbabilityDistribution: {[edgeKey: string]: number[]}
  cascadingPointsScreenPosition: {[locationID: number]: [number, number]}
  influencesDict: InfluencesDict
  influencesContrary: InfluencesDict
  timeLength: number
  allIDs: number[]
  filteredIDs: number[]
  distribution: number[]
  allLocations: ILocation[]
  allLocationsAverageDataDict: {[key: number]: [number, number][]} // [eventCount per hour or per month, averageDuration]
  allLocationsInfoDict: {[key: number]: ILocationInfo}
  locationsInfo: ILocationInfo[]
  maxCount: number
  minCount: number
  colorMapForEachLocation: {[rid: number]: string}
  clustersWithContext: ICluster[]
  clusterInfectEventsDict: {[key: number]: InfectEvent[]}
  projectionWidth: number
  projectionHeight: number
  eventsWithIDList: IEventsWithID[]
  eventsDict: {[locationID: number]: number[]}
  maxValueAcrossClusters: IMaxValueAcrossClusters
  pieChartInnerRate: number
  trees: PropagationTree[]
  selectedTreesId: number[]
  selectedAggregatePropagationTreesId: number[]
  aggregatePropagationTrees: AggregatePropagationTree[]
  colors: string[]
  locationColors: {
    AirPollution: string[],
    Congestion: string[],
    [key: string]: string[]
  }
  selectedEdge: {[edgeKey: string]: boolean}
  openedEdge: {[edgeKey: string]: boolean}
  timeWindows: ITimeWindowInfo[]
  moveLatlngDict: {[key: number]: [number, number]}
}

export const state: IProjStoreState = {
  currentDataSetConfig: {
    name: 'Congestion',
    abbrev: 'congestion2',
    mapCenter: [30.3, 120.15],
    zoom: 14,
    rangeEpsS: [0, 6],
    initParams: [0.085, 2.5, 10],
    pieChartInnerRate: 10,
    donutChartSmallRadius: 4,
    donutSectorNum: 24,
    glyphMinZoom: 14,
    secondGlyphZoom: 15,
    startTime: new Date(2018, 2, 1),
    deltaTime: 10 * 60,
    xTimeBin: 61 + 1,
    yTimeBin: 24 * 6,
    yTimeBinRate: 6,
    gridShape: [61 + 1, 24],
    totalInfluenceDuration: 16,
    pLowerBound: 0.5,
    timelineLengthRate: 300 / 18
  },
  url: 'http://localhost:8888/',
  // url: 'http://10.0.105.108:8888/',
  // url: 'http://192.168.109.128:8888/',
  embedPoints: [],
  cascadingLinks: [],
  edgeProbabilityDistribution: {},
  cascadingPointsScreenPosition: {},
  influencesDict: {},
  influencesContrary: {},
  timeLength: 350,
  allIDs: [],
  filteredIDs: [],
  distribution: [],
  allLocations: [],
  allLocationsAverageDataDict: {},
  allLocationsInfoDict: {},
  locationsInfo: [],
  maxCount: 0,
  minCount: 0,
  colorMapForEachLocation: {},
  clustersWithContext: [],
  clusterInfectEventsDict: {},
  projectionWidth: 0,
  projectionHeight: 0,
  eventsWithIDList: [],
  eventsDict: {},
  maxValueAcrossClusters: {
    maxStartTimeNumber: 0,
    maxRawEventNumber: 0,
    maxLengthNumber: 0
  },
  pieChartInnerRate: 1,
  trees: [],
  selectedTreesId: [],
  selectedAggregatePropagationTreesId: [],
  aggregatePropagationTrees: [],
  colors: ['#1f78b4', '#33a02c', '#e31a1c', '#ff7f00', '#6a3d9a', '#a6cee3', '#b2df8a', '#fb9a99',
    '#fdbf6f', '#cab2d6', '#ffff99', '#b15928'],
  // locationColors: ['#eb2f96', '#722ed1', '#13c2c2', '#1890ff', '#faad14',
  //   '#fa541c', '#b2df8a', '#fb9a99'],
  // locationColors: ['#c41d7f', '#d46b08', '#08979c', '#096dd9', '#7cb305', '#cf1322'], // congestion
  locationColors: {
    AirPollution: ['#08979c', '#096dd9', '#d46b08', '#c41d7f', '#7cb305', '#cf1322'],
    Congestion: ['#c41d7f', '#d46b08', '#08979c', '#096dd9', '#7cb305', '#cf1322']
  }, // air
  selectedEdge: {},
  openedEdge: {},
  timeWindows: [],
  moveLatlngDict: {
    10925: [30.299908621272724, 120.08775603692845],
    10948: [30.303040069845245, 120.10367424222385],
    10933: [30.289370305765935, 120.10745096787113],
    10885: [30.277758104591587, 120.08886080260257],
    10852: [30.26428634715248, 120.08211717276187],
    10895: [30.26666878983285, 120.12244476288643],
    10943: [30.277986622654797, 120.11958869186833],
    10945: [30.279273955465598, 120.12765483587899],
    10920: [30.29192764439156, 120.15441927722074],
    10865: [30.256745201603966, 120.05807767818841],
    10950: [30.28578595692056, 120.1145223057224]
  }
}
