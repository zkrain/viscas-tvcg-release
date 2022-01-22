<template>
  <div class="map-pattern-div">
    <svg width="0" height="0" class="def-svg">
      <defs>
        <filter id="blurFilter" x="-5" y="-5" width="40" height="40">
          <feGaussianBlur in="SourceAlpha" result="blurOut" stdDeviation="4"/>
          <feColorMatrix in="blurOut" type="matrix" result="colorOut"
                         values="1 1 1 0.8 0
                                  1 1 1 0.8 0
                                  1 1 1 0.8 0
                                  1 1 1 0.8 0"/>
          <feMerge>
            <feMergeNode in="colorOut" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>
    </svg>
    <div class="pmap" ref="pmap" />
    <div class="view-name">
      <div class="name-div">Spatial View</div>
    </div>
    <div  @mouseleave="showProj = false">
      <div class="projection-div" @mouseenter="showProj = true">
        <div class="projection-btn">Projection</div>
        <div class="projection-icon">
          <font-awesome-icon :icon="showProj? angleUp : angleDown" class="angle-icon"/>
        </div>
      </div>
      <transition name="slide-fade">
        <div class="projection-in-map-cont" v-show="showProj" :style="{display: 'flex'}">
          <test-projection-presentation />
        </div>
      </transition>
    </div>
    <div class="infer-box-cls">
      <infer-panel/>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import L from 'leaflet'
import 'leaflet-css'
import 'leaflet-draw/dist/leaflet.draw-src.css'
import 'leaflet-draw'
import 'animate.css'
import '@/css/transition.css'
import { faAngleDown, faAngleUp } from '@fortawesome/free-solid-svg-icons'
import * as d3 from 'd3'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { namespace } from 'vuex-class'
import { BindingHelpers } from 'vuex-class/lib/bindings'
import { indexOf, uniq, without } from 'lodash'
import bezier from 'bezier-easing'
import ProjectionInPresentation from './ProjectionInPresentation.vue'
import EventBus from '../../eventBus'
import { crossThreePoints, distance, dotThreePoints } from '@/utils/geoCompute'
import * as turf from '@turf/turf'
import InferPanel from '@/components/MapView/InferPanel.vue'
import TestProjectionPresentation from '@/components/MapView/TestProjectionPresentation.vue'

const easing: BezierEasing.EasingFunction = bezier(0, 0, 0.25, 1)

const projStore: BindingHelpers = namespace('Proj')
type LatLng = [number, number]
type LatLngLink = [LatLng, LatLng]
type LinkInfo = { locationsId: [number, number], locationsPosition: LatLngLink, rank: number}

@Component({
  components: {
    TestProjectionPresentation,
    InferPanel,
    FontAwesomeIcon,
    ProjectionInPresentation
  }
})
export default class MapView extends Vue {
  @projStore.State private allIDs!: number[]
  @projStore.State private allLocationsInfoDict!: {[key: number]: ILocationInfo}
  @projStore.State private allLocationsAverageDataDict!: {[key: number]: [number, number][]}
  @projStore.State private currentDataSetConfig!: IDataSetConfig
  @projStore.State private cascadingLinks!: number[][]
  @projStore.State private influencesDict!: InfluencesDict
  @projStore.Action private clearCascadingPattern!: () => void
  @projStore.Getter private locationDict!: {[locationID: number]: ILocation}
  @projStore.Getter private locationRange!: ILocationRange
  @projStore.Getter private edgeColorDict!: {[edgeKey: string]: [string, string]}
  @projStore.Getter private locationColorDict!: {[locationID: number]: string}
  @projStore.Action private selectEdge!: (edgeKey: string) => void
  @projStore.State private selectedEdge!: {[edgeKey: string]: boolean}
  @projStore.State private moveLatlngDict!: {[key: number]: [number, number]}
  @projStore.Action private setCascadePointsScreenPosition!: (points: {[locationID:number]: [number, number]}) => void
  private topHighlightLocationsId: number[] = []
  private bottomHighlightLocationsId: number[] = []
  private pmap!: L.Map
  private nodeG!: d3.Selection<SVGGElement, unknown, null, undefined>
  private linkG!: d3.Selection<SVGGElement, unknown, null, undefined>
  private cascadingLinksInfo: LinkInfo[] = []
  private linkWidthRange: [number, number] = [4, 8]
  private drawnItems!: L.FeatureGroup
  private defs!: d3.Selection<SVGDefsElement, unknown, null, undefined>
  private mapSvg!: d3.Selection<SVGSVGElement, unknown, null, undefined>
  private hasMarker: {[key: string]: string} = {} // {color: name}
  private showProj: boolean = false
  private maxEventCount: number = 0
  private minEventCount: number = 0

  private insertToTop (locationsId: number[]): void {
    this.topHighlightLocationsId.push(...locationsId)
    this.bottomHighlightLocationsId = without(this.bottomHighlightLocationsId, ...this.topHighlightLocationsId)
    EventBus.$emit('setInferLocationIds', this.topHighlightLocationsId)
  }

  private insertToBottom (locationsId: number[], refreshTop: boolean = true): void {
    this.bottomHighlightLocationsId.push(...locationsId)
    if (refreshTop) {
      this.topHighlightLocationsId = []
      this.fitBounds(this.bottomHighlightLocationsId)
    }
    EventBus.$emit('setInferLocationIds', this.topHighlightLocationsId)
  }

  private onProjSelect (locationsId: number[]): void {
    this.removeAll()
    this.clearCascadingPattern()
    this.topHighlightLocationsId = []
    this.bottomHighlightLocationsId = []
    this.drawnItems.clearLayers()
    this.initDrawVariable()
    this.insertToBottom(locationsId)
  }

  get locationIds (): number[] {
    let locationIds: number[] = []
    this.cascadingLinksInfo.forEach((link: LinkInfo) => {
      locationIds.push(...link.locationsId)
    })
    return uniq(locationIds)
  }

  @Watch('allLocationsInfoDict')
  private async onAllLocationsInfoDictChanged (): Promise<void> {
    this.topHighlightLocationsId = []
    this.bottomHighlightLocationsId = []
    await this.initDrawVariable()
    this.reDrawAll()
  }

  @Watch('cascadingLinks')
  private onCascadingLinksChanged (): void {
    if (this.cascadingLinks.length > 0) {
      this.drawnItems.clearLayers()
      this.initDrawVariable()
      this.reDrawAll()
      // this.fitBounds(this.topHighlightLocationsId)
    }
  }

  @Watch('topHighlightLocationsId')
  @Watch('bottomHighlightLocationsId')
  private onHighLightLocationsIdChanged (): void {
    if (typeof this.allLocationsInfoDict !== 'undefined') {
      this.reDrawAll()
    }
  }

  get angleDown () {
    return faAngleDown
  }

  get angleUp () {
    return faAngleUp
  }

  private mounted (): void {
    EventBus.$on('updateMapCenterAndZoom', this.updateMapCenterAndZoom)
    EventBus.$on('onProjSelect', this.onProjSelect)
    this.pmap = L.map(this.$refs.pmap as HTMLElement, {
      center: [30.3, 120.15],
      zoom: 12.5,
      zoomSnap: 0.1,
      zoomControl: true,
      attributionControl: false,
      touchZoom: true
    })
    // this.pmap.addEventListener('click', (evt) => {
    //   // @ts-ignore
    //   const latlng = evt.latlng
    //   let textarea = document.createElement('input')
    //   document.body.appendChild(textarea)
    //   textarea.value = latlng.lat + ',' + latlng.lng + ','
    //   console.log(textarea.value)
    //   textarea.select()
    //   document.execCommand('Copy')
    //   document.body.removeChild(textarea)
    // })
    this.drawnItems = new L.FeatureGroup()
    this.pmap.addLayer(this.drawnItems)
    const options = {
      position: 'topleft',
      draw: {
        polygon: {
          allowIntersection: false, // Restricts shapes to simple polygons
          drawError: {
            color: '#e1e100', // Color the shape will turn when intersects
            message: '<strong>Oh snap!<strong> you can\'t draw that!' // Message that will show when intersect
          },
          shapeOptions: {
            color: '#97009c'
          }
        },
        // disable toolbar item by setting it to false
        polyline: false,
        circle: false,
        rectangle: false,
        marker: false
      },
      edit: {
        featureGroup: this.drawnItems, // REQUIRED!!
        remove: true
      }
    }
    // @ts-ignore
    const drawControl = new L.Control.Draw(options)
    this.pmap.addControl(drawControl)
    // @ts-ignore
    this.pmap.on(L.Draw.Event.CREATED, (e: L.ResizeEvent) => {
      const layer: L.Layer = e.layer
      const polygonPoints: L.LatLng[] = e.layer.editing.latlngs[0][0]
      polygonPoints.push(polygonPoints[0])
      this.polygonSelect(polygonPoints)
      this.drawnItems.addLayer(layer)
    })
    // const tileUrl: string = 'https://api.mapbox.com/styles/v1/zikunrain/ck6d7y0t421y71io28j9no6ka/tiles/256/{z}/{x}/{y}@2x?access_token=pk.eyJ1IjoiemlrdW5yYWluIiwiYSI6ImNqeWE2dXJ1djBibmIzY21mMWl5MDljc2wifQ.NMe8T_yYFKIsraDJV4tIPw'
    // const tileUrl: string = 'http://a.tile.openstreetmap.org/{z}/{x}/{y}.png'
    // const tileUrl: string = 'http://mt.google.com/vt/lyrs=m&x={x}&y={y}&z={z}'
    const tileUrl: string = 'https://api.mapbox.com/styles/v1/zikunrain/ckev03ci403591as420c53a2x/tiles/256/{z}/{x}/{y}@2x?access_token=pk.eyJ1IjoiemlrdW5yYWluIiwiYSI6ImNqeWE2dXJ1djBibmIzY21mMWl5MDljc2wifQ.NMe8T_yYFKIsraDJV4tIPw'

    L.tileLayer(tileUrl, {
      tileSize: 512,
      zoomOffset: -1
    }).addTo(this.pmap)
    this.pmap.on('zoomstart', this.zoomStart)
    this.pmap.on('zoomend', this.zoomEnd)
    this.initSVG()
    if (typeof this.allLocationsInfoDict !== 'undefined') {
      this.initDrawVariable()
      this.fitBoundDataSet()
    }
  }

  private zoomStart (): void {
    this.removeAll()
  }

  private zoomEnd (): void {
    this.reDrawAll()
  }

  private polygonSelect (polygonPoints: L.LatLng[]): void {
    this.insertToBottom(this.topHighlightLocationsId, false)
    this.topHighlightLocationsId = []
    this.removeAll()
    this.cascadingLinksInfo = []
    this.clearCascadingPattern()
    this.removeAll()
    const selectLocationId: number[] = []
    const points: number[][] = []
    polygonPoints.forEach((pt: L.LatLng) => {
      points.push([pt.lat, pt.lng])
    })
    const tmpPolygonPoints: number[][][] = []
    tmpPolygonPoints.push(points)
    const polygon = turf.polygon(tmpPolygonPoints)
    let locationsId: number[] = []
    if (this.bottomHighlightLocationsId.length > 0) {
      locationsId = this.bottomHighlightLocationsId
    } else {
      locationsId = this.allIDs
    }
    locationsId.forEach((locationId: number) => {
      const location = this.allLocationsInfoDict[locationId].location
      const pt = turf.point([location.lat, location.lng])
      if (turf.booleanContains(polygon, pt)) {
        selectLocationId.push(locationId)
      }
    })
    this.insertToTop(selectLocationId)
    const leftBottomPt: LatLng = [90, 180]
    const rightTopPt: LatLng = [-90, -180]
    polygonPoints.forEach((point: L.LatLng) => {
      leftBottomPt[0] = Math.min(leftBottomPt[0], point.lat)
      leftBottomPt[1] = Math.min(leftBottomPt[1], point.lng)
      rightTopPt[0] = Math.max(rightTopPt[0], point.lat)
      rightTopPt[1] = Math.max(rightTopPt[1], point.lng)
    })
    this.pmap.fitBounds([leftBottomPt, rightTopPt])
  }

  private fitBoundDataSet (): void {
    this.updateMapCenterAndZoom(this.currentDataSetConfig.mapCenter, this.currentDataSetConfig.zoom)
  }

  private fitBounds (locationsId: number[]): void {
    console.log('fitBounds')
    const leftBottomPt: LatLng = [90, 180]
    const rightTopPt: LatLng = [-90, -180]
    locationsId.forEach((locationId: number) => {
      const location: ILocation = this.allLocationsInfoDict[locationId].location
      leftBottomPt[0] = Math.min(leftBottomPt[0], location.lat)
      leftBottomPt[1] = Math.min(leftBottomPt[1], location.lng)
      rightTopPt[0] = Math.max(rightTopPt[0], location.lat)
      rightTopPt[1] = Math.max(rightTopPt[1], location.lng)
    })
    this.pmap.fitBounds([[leftBottomPt[0] + 0.01, leftBottomPt[1] + 0.02], [rightTopPt[0] - 0.02, rightTopPt[1] - 0.02]])
    // this.pmap.fitBounds([leftBottomPt, rightTopPt])
  }

  private initOneSVG (svgName: string): d3.Selection<SVGSVGElement, unknown, null, undefined> {
    const width: number = this.pmap.getSize().x
    const height: number = this.pmap.getSize().y
    return d3.select(this.pmap.getPanes().overlayPane)
      .append('svg')
      .attr('class', svgName)
      .attr('width', width)
      .attr('height', height)
      .attr('overflow', 'visible')
  }

  private initSVG (): void {
    this.mapSvg = this.initOneSVG('map-svg')
    this.defs = this.mapSvg.append('defs')
    this.linkG = this.mapSvg.append('g').attr('class', 'links-g')
    this.nodeG = this.mapSvg.append('g').attr('class', 'nodes-g')
  }

  private initDrawVariable (): void {
    this.initCascadingLinksInfo()
    this.initNodes()
  }

  private initDefs (): void {
    this.defs.selectAll('marker').remove()
    this.hasMarker = {}
    this.addArrowMarkerMarker('#516172', 'gray-arrow')
    for (let key in this.edgeColorDict) {
      const keysNum = Object.keys(this.hasMarker).length
      const color: string = this.edgeColorDict[key][1]
      const id: string = `color_${keysNum}`
      this.addArrowMarkerMarker(color, id)
    }
  }

  private addArrowMarkerMarker (color: string, id: string) {
    const markerWidth: number = 4
    this.defs.append('marker').attr('id', id)
      .attr('markerWidth', markerWidth)
      .attr('markerHeight', markerWidth / 3 * 2)
      .attr('refX', 2.5)
      .attr('refY', markerWidth / 3)
      .attr('orient', 'auto')
      .attr('markerUnits', 'strokeWidth')
      .append('path')
      .attr('d', `M0,0 L0,${markerWidth / 3 * 2} L${markerWidth},${markerWidth / 3} Z`)
      .attr('fill', color)
    this.hasMarker[color] = id
  }

  private initCascadingLinksInfo (): void {
    this.cascadingLinksInfo = []
    this.cascadingLinks.forEach((link: number[]) => {
      let startPt: ILocation = this.locationDict[link[0]]
      if (startPt.rid in this.moveLatlngDict) {
        startPt.lat = this.moveLatlngDict[startPt.rid][0]
        startPt.lng = this.moveLatlngDict[startPt.rid][1]
      }
      let endPt: ILocation = this.locationDict[link[1]]
      if (endPt.rid in this.moveLatlngDict) {
        endPt.lat = this.moveLatlngDict[endPt.rid][0]
        endPt.lng = this.moveLatlngDict[endPt.rid][1]
      }
      const linkPosition: LatLngLink = [[startPt.lat, startPt.lng], [endPt.lat, endPt.lng]]
      this.cascadingLinksInfo.push({ locationsId: [link[0], link[1]], locationsPosition: linkPosition, rank: link[2] })
    })
    console.log('cascading links info')
    console.log(this.cascadingLinksInfo)
  }

  private initNodes (): void {
    let locationIds: number[] = this.locationIds
    this.insertToBottom(without(this.topHighlightLocationsId, ...locationIds), false)
    this.topHighlightLocationsId = locationIds
    const screenPointsPosition: {[locationID:number]: [number, number]} = {}
    locationIds.forEach((locationId: number) => {
      const pt: ILocation = this.locationDict[locationId]
      const mapPt: L.Point = this.pmap.latLngToLayerPoint({ lat: pt.lat, lng: pt.lng })
      screenPointsPosition[locationId] = [mapPt.x, mapPt.y]
    })
    this.setCascadePointsScreenPosition(screenPointsPosition)
  }

  private removeAll (): void {
    this.nodeG.selectAll('*').remove()
    this.linkG.selectAll('*').remove()
  }

  private reDrawAll (): void {
    this.removeAll()
    if (this.bottomHighlightLocationsId.length > 0 || this.topHighlightLocationsId.length > 0) { // 投影图中选了或者地图中选了的情况
      this.drawPoints('points',
        without(this.allIDs, ...this.topHighlightLocationsId, ...this.bottomHighlightLocationsId),
        false, [3, 3]) // 未选择的点
      if (this.pmap.getZoom() >= this.currentDataSetConfig.glyphMinZoom) {
        this.drawDonutChart('bottomHighlight', this.bottomHighlightLocationsId) // 投影图选了的点, 显示为donut chart
      } else {
        this.drawPoints('bottomHighlight', this.bottomHighlightLocationsId, true,
          [6, 10], true) // 投影图选了的点, 显示为点
      }
      this.drawDonutChart('highlight', this.topHighlightLocationsId) // 地图中选了的点
      this.drawLinks()
    } else { // 什么也不选的情况
      if (this.pmap.getZoom() >= this.currentDataSetConfig.secondGlyphZoom) {
        const bounds = this.pmap.getBounds()
        const drawIds: number[] = []
        this.allIDs.forEach((id: number) => {
          const pt = this.allLocationsInfoDict[id].location
          if (bounds.contains({ lat: pt.lat, lng: pt.lng })) {
            drawIds.push(id)
          }
        })
        this.drawDonutChart('points', drawIds)
      } else {
        this.drawPoints('points', this.allIDs, true)
      }
    }
  }

  private drawLinks (): void {
    this.linkG.selectAll('.link').remove()
    this.initDefs()
    const k: number = (this.linkWidthRange[1] - this.linkWidthRange[0]) / this.cascadingLinksInfo.length
    const linkNum: number = this.cascadingLinksInfo.length
    this.linkG.selectAll('.link')
      .data(this.cascadingLinksInfo)
      .enter()
      .append('path')
      .attr('class', (t: LinkInfo) => 'link' + t.locationsId[0] + '-' + t.locationsId[1])
      .attr('fill', 'none')
      .attr('stroke-width', (t: LinkInfo) => this.linkWidthRange[0] + k * (linkNum - t.rank))
      .attr('stroke', '#516172')
      .attr('filter', (t: LinkInfo) => {
        const key: string = t.locationsId[0] + ',' + t.locationsId[1]
        if (key in this.selectedEdge) {
          return 'url(#blurFilter)'
        }
        return ''
      })
      .attr('marker-end', 'url(#gray-arrow)')
      .attr('d', (link: LinkInfo) => {
        return this.latLngLinkToPath(link.locationsPosition, this.pmap)
      })
    this.linkG.selectAll('.link_mouse')
      .data(this.cascadingLinksInfo)
      .enter()
      .append('path')
      .attr('class', 'link_mouse')
      .attr('fill', 'none')
      .attr('stroke-width', (t: LinkInfo) => this.linkWidthRange[0] + k * (linkNum - t.rank) + 10)
      .attr('stroke', '#516172')
      .attr('opacity', 0)
      .attr('marker-end', 'url(#gray-arrow)')
      .attr('cursor', 'pointer')
      .attr('d', (link: LinkInfo) => {
        return this.latLngLinkToPath(link.locationsPosition, this.pmap)
      })
      .on('click', (link: LinkInfo) => {
        const linkKey: string = `${link.locationsId[0]},${link.locationsId[1]}`
        this.selectEdge(linkKey)
      })
  }

  @Watch('selectedEdge', { deep: true })
  private onSelectedEdgeChanged () {
    d3.selectAll('path').attr('filter', '')
    for (let selectedEdgeKey in this.selectedEdge) {
      const className: string = 'link' + selectedEdgeKey.replace(',', '-')
      d3.select('.' + className).attr('filter', 'url(#blurFilter)')
    }
  }

  private drawPoints (className: string, locationsId: number[], highlight: boolean = true,
    radiusRange: [number, number] = [3, 6], encodeWithOpacity: boolean = false): void {
    // if radius == -1, encode number of total events with radius (range： 3-6)
    if (typeof locationsId === 'undefined') {
      return
    }
    this.initMaxMinEventCount(locationsId)
    const pointClassName: string = `${className}_circle`
    const tLocationsInfo: ILocationInfo[] = []
    locationsId.forEach((locId: number) => {
      tLocationsInfo.push(this.allLocationsInfoDict[locId])
    })
    if (tLocationsInfo.length === 0) {
      return
    }
    const rs: number[] = new Array(tLocationsInfo.length).fill(3)
    let maxVal = 0
    let minVal = 999999999
    tLocationsInfo.forEach((t) => {
      maxVal = Math.max(t.totalEventCount, maxVal)
      minVal = Math.min(t.totalEventCount, minVal)
    })
    const k = 1 / Math.max(maxVal - minVal, 1)
    tLocationsInfo.forEach((t, i) => {
      rs[i] = radiusRange[0] + (radiusRange[1] - radiusRange[0]) * k * (t.totalEventCount - minVal)
    })
    this.nodeG.selectAll('.' + pointClassName)
      .data(tLocationsInfo)
      .enter()
      .append('circle')
      .attr('class', pointClassName)
      .attr('id', t => t.location.rid)
      .attr('r', (t, i) => rs[i])
      .attr('fill', '#777')
      .attr('stroke', 'white')
      .attr('stroke-width', 1)
      .attr('opacity', (t: ILocationInfo) => {
        if (highlight && encodeWithOpacity) {
          return this.getOpacity(t.location.rid)
        } else if (highlight && !encodeWithOpacity) {
          return 1
        }
        return 0.1
      })
      .attr('transform', (locationInfo: ILocationInfo) => {
        const location: ILocation = locationInfo.location
        if (location.rid in this.moveLatlngDict) {
          location.lat = this.moveLatlngDict[location.rid][0]
          location.lng = this.moveLatlngDict[location.rid][1]
        }
        const pt: L.Point = this.pmap.latLngToLayerPoint(
          { lat: location.lat, lng: location.lng })
        return `translate(${pt.x}, ${pt.y})`
      })
      // .on('click', (e) => {
      //   console.log(e.location.rid)
      // })
      .attr('cursor', (t: ILocationInfo) => {
        if (highlight && encodeWithOpacity) {
          return this.getOpacity(t.location.rid)
        } else if (highlight && !encodeWithOpacity) {
          return 'pointer'
        }
        return 'inherit'
      })
  }

  private initMaxMinEventCount (locationsId: number[]): void {
    const totalCounts: number[] = locationsId.map((id: number) => {
      return this.allLocationsInfoDict[id].totalEventCount
    })
    this.maxEventCount = Math.max(...totalCounts)
    this.minEventCount = Math.min(...totalCounts)
  }

  private getOpacity (thisId: number): number {
    const val: number = (this.allLocationsInfoDict[thisId].totalEventCount - this.minEventCount) / (this.maxEventCount - this.minEventCount)
    return 0.2 + 0.8 * val
  }

  private drawDonutChart (className: string, locationsId: number[]): void {
    const pieChartGCls: string = `${className}_g`
    const arcCls: string = `${className}_arc`
    const outerCircleCls: string = `${className}_outer_circle`
    const midCircleRadius: number = 24
    const outerCircleRadius: number = 18
    const maxInnerCircleRadius: number = 16
    const tLocationsInfo: ILocationInfo[] = []

    let maxEventCount: number = 0
    let minEventCount: number = 999999
    this.allIDs.forEach((key: number) => {
      maxEventCount = Math.max(this.allLocationsInfoDict[key].totalEventCount, maxEventCount)
      minEventCount = Math.min(this.allLocationsInfoDict[key].totalEventCount, minEventCount)
    })
    locationsId.forEach((key: number) => {
      if (indexOf(this.topHighlightLocationsId, key) > 0) {
        tLocationsInfo.push(this.allLocationsInfoDict[key])
      } else {
        tLocationsInfo.unshift(this.allLocationsInfoDict[key])
      }
    })
    if (tLocationsInfo.length === 0) {
      return
    }
    this.nodeG.selectAll('.' + pieChartGCls)
      .data(tLocationsInfo)
      .enter()
      .append('g')
      .attr('class', pieChartGCls)
      .attr('id', (loc: ILocationInfo) => loc.location.rid)
      // .on('click', (e) => {
      //   console.log('rid:' + e.location.rid)
      // })
      .attr('cursor', 'pointer')
      .call((g: d3.Selection<SVGGElement, ILocationInfo, SVGGElement, unknown>) => {
        g.append('circle')
          .attr('class', outerCircleCls)
          .attr('fill', (locationInfo: ILocationInfo) => {
            const rid: number = locationInfo.location.rid
            if (rid in this.locationColorDict) {
              return this.locationColorDict[rid]
            }
            return '#999'
          })
          .attr('r', (locationInfo: ILocationInfo): number => {
            const normalizedTotalCount: number = (locationInfo.totalEventCount - minEventCount) /
              (maxEventCount - minEventCount)
            return maxInnerCircleRadius * (0.4 + 0.6 * normalizedTotalCount)
          })
          .attr('transform', (locationInfo: ILocationInfo) => {
            const location = locationInfo.location
            if (location.rid in this.moveLatlngDict) {
              location.lat = this.moveLatlngDict[location.rid][0]
              location.lng = this.moveLatlngDict[location.rid][1]
            }
            const pt: L.Point = this.pmap.latLngToLayerPoint({ lat: location.lat, lng: location.lng })
            return `translate(${pt.x},${pt.y})`
          })
      })
    const arc: d3.Arc<any, any> = d3.arc().innerRadius(outerCircleRadius).outerRadius(midCircleRadius)
    tLocationsInfo.forEach((locationInfo: ILocationInfo, i: number) => {
      let color: string = '#777'
      if (this.locationColorDict.hasOwnProperty(locationInfo.location.rid)) {
        color = this.locationColorDict[locationInfo.location.rid]
      }
      const location = locationInfo.location
      if (location.rid in this.moveLatlngDict) {
        location.lat = this.moveLatlngDict[location.rid][0]
        location.lng = this.moveLatlngDict[location.rid][1]
      }
      const pt: L.Point = this.pmap.latLngToLayerPoint({ lat: location.lat, lng: location.lng })
      const rid: number = locationInfo.location.rid
      const aveList: [number, number][] = this.allLocationsAverageDataDict[rid]
      let minV: number = 99999999
      let maxV: number = 0
      this.allLocationsAverageDataDict[rid].forEach((data: [number, number]) => {
        minV = Math.min(data[0], minV)
        maxV = Math.max(data[0], maxV)
      })
      this.nodeG.append('g')
        .attr('transform', () => `translate(${pt.x},${pt.y})`)
        .attr('fill', color)
        .selectAll('.' + arcCls + i)
        .data(aveList)
        .enter()
        .append('path')
        .attr('opacity', (data: [number, number]) => {
          const val = (data[0] - minV) / (maxV - minV)
          return 0.1 + 0.8 * val * val
        })
        .attr('d', (data: [number, number], key: number) => {
          return arc({ startAngle: 2 * Math.PI * key / this.currentDataSetConfig.donutSectorNum,
            endAngle: 2 * Math.PI * (key + 1) / this.currentDataSetConfig.donutSectorNum })
        })
    })
  }

  private updateMapCenterAndZoom (center: [number, number], zoom: number): void {
    this.pmap.setView({ lat: center[0], lng: center[1] }, zoom)
  }

  private latLngLinkToPath (link: LatLngLink, map: L.Map, padding: number = 24): string {
    let pt0: L.Point = map.latLngToLayerPoint({ lat: link[0][0], lng: link[0][1] })
    let pt1: L.Point = map.latLngToLayerPoint({ lat: link[1][0], lng: link[1][1] })
    const line1 = distance([pt0.x, pt0.y], [pt1.x, pt1.y])
    const ang: number = Math.atan2(pt1.y - pt0.y, pt1.x - pt0.x)
    let controlRate: number = 0.35 // turn right when it is positive, turn left when it is negative
    let moveRate: number = 0.7 // turn right when it is positive, turn left when it is negative
    let left: number = 0
    let right: number = 0
    const locationIds: number[] = this.locationIds
    locationIds.forEach((id: number) => {
      const pt: L.Point = map.latLngToLayerPoint({ lat: this.locationDict[id].lat, lng: this.locationDict[id].lng })
      if ((Math.abs(pt0.x - pt.x) < 1e-6 && Math.abs(pt0.y - pt.y) < 1e-6) ||
        (Math.abs(pt1.x - pt.x) < 1e-6 && Math.abs(pt1.y - pt.y) < 1e-6)) {
        // do nothing
      } else {
        const cross = crossThreePoints([pt0.x, pt0.y], [pt1.x, pt1.y], [pt.x, pt.y])
        const proj = Math.abs(dotThreePoints([pt0.x, pt0.y], [pt1.x, pt1.y], [pt.x, pt.y]) /
          distance([pt0.x, pt0.y], [pt1.x, pt1.y]))
        if (cross > 0 && line1 >= proj) {
          left++
        } else if (cross < 0 && line1 >= proj) {
          right++
        }
      }
    })
    if (left > right) {
      controlRate *= -1
      moveRate *= -1
    }
    pt0 = new L.Point(pt0.x - moveRate * padding * Math.sin(ang), pt0.y + moveRate * padding * Math.cos(ang))
    pt1 = new L.Point(pt1.x - moveRate * (padding + 15) * Math.sin(ang), pt1.y + moveRate * (padding + 15) * Math.cos(ang))
    const dx: number = padding * Math.cos(ang)
    const dy: number = padding * Math.sin(ang)
    const startPt: [number, number] = [pt0.x + dx, pt0.y + dy]
    const endPt: [number, number] = [pt1.x - dx, pt1.y - dy]
    const lineLength: number = distance(startPt, endPt)
    let controlPt: [number, number] = [0.5 * (startPt[0] + endPt[0]) - controlRate * lineLength * Math.sin(ang),
      0.5 * (startPt[1] + endPt[1]) + controlRate * lineLength * Math.cos(ang)]
    if (line1 <= 150) {
      controlPt = [0.5 * (startPt[0] + endPt[0]) - 0.3 * controlRate * lineLength * Math.sin(ang),
        0.5 * (startPt[1] + endPt[1]) + 0.3 * controlRate * lineLength * Math.cos(ang)]
    }
    return `M${startPt[0]},${startPt[1]} Q${controlPt[0]},${controlPt[1]} ${endPt[0]},${endPt[1]}`
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.map-pattern-div {
  position: relative;
  width: 100%;
  height: 100%;

  .pmap {
    position: relative;
    width: 100%;
    height: 100%;
  }

  .view-name {
    position: absolute;
    top: 0px;
    left: 0px;
    height: 30px;
    width: 160px;
    background:rgba(8,121,235,1);
    box-shadow:0px 2px 4px 0px rgba(32,84,159,0.4);

    .name-div {
      position: relative;
      top: 6px;
      left: 32px;
      width:96px;
      height:18px;
      font-size:14px;
      font-family:Helvetica-Bold,Helvetica;
      font-weight:bold;
      color:rgba(255,255,255,1);
      line-height:17px;
      letter-spacing:1px;
    }
  }

  .def-svg {
    position: absolute;
    pointer-events: none;
    opacity: 0;
    left: 0;
    top: 0;
  }

  .map-view-btn {
    position: absolute;
    left: calc(100% - 100px);
    top: 2px;
    width: 100px;
    height: 24px;
    background-color: #fff;
    border-radius: 2px;
    border: 1px solid #eee;
    line-height: 24px;
    user-select: none;
    cursor: pointer;
  }

  .infer-btn {
    left: calc(100% - 220px);
  }

  .projection-div {
    position: absolute;
    top: 4px;
    left: calc(100% - 150px);
    width: 130px;
    height: 26px;
    box-sizing: border-box;
    background-color: white;
    border: 1px solid #eee;
    border-radius: 4px;
    transition: all 300ms;
    user-select: none;
    z-index: 2;
    display: flex;
    align-items: center;
    vertical-align: middle;
    &:hover {
      color: #3399FF;
    }

    .projection-btn {
      width: 100px;
      float: left;
      box-sizing: border-box;
      text-align: center;
      height: fit-content;
    }

    .projection-icon {
      float: left;
      width: 26px;
      height: 100%;
      display: table-cell;
      vertical-align: middle;
      .angle-icon {
        width: 100%;
        height: 100%;
      }
    }
  }

  .projection-in-map-cont {
    position: absolute;
    left: calc(100% - 425px);
    top: 15px;
    height: 420px;
    width: 420px;
    background-color: #fff;
    display: flex;
    border: 1px solid #eee;
  }

  .infer-box-cls {
    position: absolute;
    top: 4px;
    left: calc(100% - 270px);
  }
}

// .leaflet-top {
//   top: 80px !important;
// }
</style>
