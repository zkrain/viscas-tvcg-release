function distance (pt0: [number, number], pt1: [number, number]): number {
  return Math.sqrt((pt1[1] - pt0[1]) * (pt1[1] - pt0[1]) + (pt1[0] - pt0[0]) * (pt1[0] - pt0[0]))
}

// 返回角012的角度[0,pi]，注意顺序
function angleThreePoints (pt0: [number, number], pt1: [number, number], pt2: [number, number]): number {
  const vec10: [number, number] = [pt0[0] - pt1[0], pt0[1] - pt1[1]]
  const vec12: [number, number] = [pt2[0] - pt1[0], pt2[1] - pt1[1]]
  const cross: number = vec10[0] * vec12[0] + vec10[1] * vec12[1]
  const dd: number = Math.sqrt((vec12[0] * vec12[0] + vec12[1] * vec12[1]) *
    (vec10[0] * vec10[0] + vec10[1] * vec10[1]))
  if (dd < 1e-6) {
    return -1
  }
  return Math.acos(cross / dd)
}

function crossThreePoints (pt0: [number, number], pt1: [number, number], pt2: [number, number]) {
  const vec1 = [pt1[0] - pt0[0], pt1[1] - pt0[1]]
  const vec2 = [pt2[0] - pt0[0], pt2[1] - pt0[1]]
  return vec1[0] * vec2[1] - vec1[1] * vec2[0]
}

function dotThreePoints (pt0: [number, number], pt1: [number, number], pt2: [number, number]) {
  const vec1 = [pt1[0] - pt0[0], pt1[1] - pt0[1]]
  const vec2 = [pt2[0] - pt0[0], pt2[1] - pt0[1]]
  return vec1[0] * vec2[0] + vec1[1] * vec2[1]
}

export {
  distance,
  angleThreePoints,
  crossThreePoints,
  dotThreePoints
}
