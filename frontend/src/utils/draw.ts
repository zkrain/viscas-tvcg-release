import L, { LineUtil } from 'leaflet'

function sectorPathD (point: L.Point, r: number, angRate: number, clockWise: boolean, fill: boolean): string {
  let x = point.x
  let y = point.y
  let endX = x + r * Math.sin(2 * Math.PI * (0.5 - angRate))
  let endY = y + r * Math.cos(2 * Math.PI * (0.5 - angRate))
  if (angRate > 1 - 1e-6) {
    angRate = 1 - 1e-6
  }
  if (clockWise && fill) {
    if (angRate < 0.5) {
      return `M${x},${y - r}, A${r},${r} 0 0,1 ${endX},${endY} L${x},${y} Z`
    } else {
      return `M${x},${y - r}, A${r},${r} 0 1,1 ${endX},${endY} L${x},${y} Z`
    }
  }
  if ((!clockWise) && fill) {
    if (angRate < 0.5) {
      return `M${x},${y - r}, A${r},${r} 0 1,0 ${endX},${endY} L${x},${y} Z`
    } else {
      return `M${x},${y - r}, A${r},${r} 0 0,0 ${endX},${endY} L${x},${y} Z`
    }
  } else if (clockWise && (!fill)) {
    if (angRate < 0.5) {
      return `M${x},${y - r}, A${r},${r} 0 0,1 ${endX},${endY}`
    } else {
      return `M${x},${y - r}, A${r},${r} 0 1,1 ${endX},${endY}`
    }
  } else if ((!clockWise) && (!fill)) {
    if (angRate < 0.5) {
      return `M${x},${y - r}, A${r},${r} 0 1,0 ${endX},${endY}`
    } else {
      return `M${x},${y - r}, A${r},${r} 0 0,0 ${endX},${endY}`
    }
  }
  return ``
}

function arcAngleToPath (angle: number, center: [number, number], r: number): string {
  if (angle > 2 * Math.PI - 1e-6) {
    angle = angle - 1e-6
  }
  let x = center[0]
  let y = center[1]
  let endX = x + r * Math.sin(Math.PI - angle)
  let endY = y + r * Math.cos(Math.PI - angle)
  if (angle < Math.PI) {
    return `M${x},${y - r}, A${r},${r} 0 0,1 ${endX},${endY}`
  } else {
    return `M${x},${y - r}, A${r},${r} 0 1,1 ${endX},${endY}`
  }
}

function arcAngleRangePathD (angleRange: [number, number], center: [number, number], r: number) {
  const eps: number = 1e-6
  let startAngle: number = Math.min(angleRange[0], 2 * Math.PI - eps)
  startAngle = Math.max(startAngle, 0)
  let endAngle: number = Math.max(angleRange[1], eps)
  endAngle = Math.min(endAngle, 2 * Math.PI - eps)
  let startX: number = center[0] + r * Math.sin(Math.PI - startAngle)
  const startY: number = center[1] + r * Math.cos(Math.PI - startAngle)
  const endX: number = center[0] + r * Math.sin(Math.PI - endAngle)
  const endY: number = center[1] + r * Math.cos(Math.PI - endAngle)
  if (Math.abs(startX - endX) < 1) {
    startX = startX + (endAngle > Math.PI ? 1 : -1)
  }
  return `M${startX},${startY}, A${r},${r} 0 ${endAngle - startAngle < Math.PI ? 0 : 1},1 ${endX},${endY}`
}

function hexToRgb (hex: string): {[key: string]: number} | null {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result ? {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16)
  } : null
}

function getBoxPlotInfo (data: number[]): number[] {
  if (data.length === 0) {
    return []
  }
  const dataS: number[] = data.sort((a, b) => a - b)
  const length: number = dataS.length
  return [dataS[0], dataS[Math.round(0.25 * length)], dataS[Math.round(0.5 * length)],
    dataS[Math.round(0.75 * length)], dataS[length - 1]]
}

function getDateFromNumber (startTime: Date, deltaTimeSecond: number, n: number): string {
  const d: Date = new Date(startTime.getTime() + 1000 * deltaTimeSecond * n)
  const monthStr: string = '' + (d.getMonth() + 1)
  const dayStr: string = '' + d.getDate()
  const hourStr: string = ('0' + d.getHours()).substr(-2)
  const minuteStr: string = ('0' + d.getMinutes()).substr(-2)
  return `${monthStr}/${dayStr} ${hourStr}:${minuteStr}`
}

export {
  sectorPathD,
  arcAngleToPath,
  arcAngleRangePathD,
  hexToRgb,
  getBoxPlotInfo,
  getDateFromNumber
}
