interface IDatasetDescription {
  eventsFile: string
  locationFile: string
  embeddingCoordinatesFile: string
  mapConfig: string // todo
}

interface RoadShape {
  Len: number
  Shapes: string
  Lv: number
  ID: number
}

interface IEventsWithID {
  locationID: number
  events: number[]
}

type CongestEvents = number[]

interface Road {
  roadShape: RoadShape,
  events: CongestEvents
  id: number
}

interface RoadWithoutShape {
  events: CongestEvents
  id: number
}

interface GetCongestedRoadsResponse {
  Code: number
  RoadShapes: RoadShape[]
}

interface ProjPoint {
  rid: number
  left: number
  top: number
}

interface Influence {
  StartTimeI: number
  EndTimeI: number
  StartTimeJ: number
  EndTimeJ: number
  P: number
  W: number
  // Strength: number
}

interface InfluencesWithID {
  influences: Influence[]
  locationIDi: number
  locationIDj: number
}

interface InfluencesDict {
  [locationIDi: number]: {
    [locationIDj: number]: Influence[]
  }
}

interface PotentialCause {
  LocationID: number
  StartTime: number
  Duration: number
  P: number
}

interface FullInfluenceInfo {
  LocationIDi: number
  StartTimeI: number
  DurationI: number
  LocationIDj: number
  StartTimeJ: number
  DurationJ: number
  P: number
  OtherPotentialCauses: PotentialCause[]
}

interface SimplifiedInfluenceInfo {
  sti: number
  reti: number // relative ent time of I
  rstj: number // relative start time of J
  retj: number // relative ent time of J
  p: number // probability
  duration: number
  csoOrderedByLocation: PotentialCause[]
  opcs: PotentialCause[]
  pso: [number, number][] // probabilities of other potential influences
  durationLinks: [number, number, number][] // [start, end, boolean]
  durationLinksAlignByRstj: [number, number, number][]
  rankOfDuration: number
  finalP: number
}

interface PropagationTree {
  TreeStartTime: number
  TreeEndTime: number
  FullInfluenceInfoList: FullInfluenceInfo[]
  durationDict: {[locationID: number]: {
    [startTime: number]: number
  }}
  averageDuration: number
  durations: [number, number][] // [duration, id][]
  involvedEdges: [number, number][]
  involvedEdgeDict: {[edgeKey: string]: SimplifiedInfluenceInfo}
  id: number
  parentAid: number
}

// interface AggregateEdge {
//   influenceInfoList: SimplifiedInfluenceInfo[]
// }

interface AggregatePropagationTree {
  trees: PropagationTree[]
  id: number
  selected: boolean
  height: number
  involvedEdgeDictInAggregation: {[edgeKey: string]: SimplifiedInfluenceInfo[]}
  involvedEdges: [number, number][]
}

// interface CompactEventsTimeSegment {
//  nEvent: number
//  nTimeRange: number
//  orderedTimeRangeLength: number[]
// }

// type CompactEvents = CompactEventsTimeSegment[]

interface ListWithVariance {
  List: number[]
  Variance: number[]
}

interface ClusterInformation {
  // StartTimeNumberInformation: ListWithVariance
  // RawEventNumberInformation: ListWithVariance
  // LengthInformation: ListWithVariance

  StartTimeNumberMatrix: number[][]
  RawEventNumberMatrix: number[][]
  LengthMatrix: number[][]
}

interface ICluster {
  color: string
  rIDs: number[]
  hull: string
  hullPoints: [number, number][]
  clusterEventCount: number[]
  clusterAverageEventCount: number
  clusterAverageInfectEventCount: number
  clusterInformation: ClusterInformation
  spatialDistribution: number[]
  spatialRelativeCenter: {rx: number, ry: number}
}

interface ILocation {
  rid: number
  lat: number
  lng: number
}

interface ILocationStyle {
  angRate: number
  normalizedAngRate: number
  aveDuration: number
}

interface ILocationInfo {
  location: ILocation
  totalEventCount: number
  locationStyle: ILocationStyle
}

interface IDataSetConfig {
  name: string
  abbrev: string
  mapCenter: [number, number]
  zoom: number
  rangeEpsS: [number, number]
  initParams: [number, number, number]
  pieChartInnerRate: number // 地图中扇形图中小扇形的缩放比例
  donutChartSmallRadius: number // donut chart中小圆的最大半径
  donutSectorNum: number
  glyphMinZoom: number
  secondGlyphZoom: number
  startTime: Date
  deltaTime: number
  xTimeBin: number
  yTimeBin: number
  yTimeBinRate: number
  gridShape: [number, number]
  totalInfluenceDuration: number
  pLowerBound: number
  timelineLengthRate: number
}

interface IMaxValueAcrossClusters {
  maxStartTimeNumber: number
  maxRawEventNumber: number
  maxLengthNumber: number
}

interface SpatialDistributionOfClusters {
  MaxDistance: number
  MinDistance: number
  DistributionList: number[][]
}

interface ILocationRange {
  minLat: number
  maxLat: number
  minLng: number
  maxLng: number
}

interface InfectEvent {
  startTime: number
  durationTime: number
  endTime: number
}

interface ITimeWindowInfo {
  timeRange: [number, number]
  tree: PropagationTree
}
