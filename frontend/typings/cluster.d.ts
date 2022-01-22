interface ICircleInfo {
  readonly clusterId: number
  rids: number[]
  color: string
  arcPath: string
  boxPlotPath: string
  center: [number, number]
  radius: number
  hullPoints: [number, number][]
  drawCenterPosition: [number, number]
}

interface IMapInfo {
  readonly clusterId: number
  mapRelativePosition: [number, number] // [left, top]
  zIndex: number
  locationIds: number[]
  locationsInfo: ILocationInfo[]
}
