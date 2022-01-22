function stdbscan(D: ILocation[], EpsT: number, EpsS: number, MinPts: number, deltaE: number, embedPointDict: {[rid: number]: ProjPoint}): { clusters: number[][], noise: number[] } {

  const clusters: number[][] = []
  let currentCluster: number[] = []
  let currentClusterLabel = -1
  const noise: number[] = [] // lable -1
  const clusterMap: {[rid: number]: number} = {}
  const Q: ILocation[] = []

  D.forEach((d: ILocation, index: number) => {
    if (!(d.rid in clusterMap)) {
      const X: ILocation[] = Retrieve_Neighbors(d, D, EpsT, EpsS, embedPointDict)
      if (X.length < MinPts) {
        clusterMap[d.rid] = -1
        noise.push(d.rid)
      } else {
        currentClusterLabel += 1
        if (currentCluster.length > 0) {
          clusters.push(currentCluster)
        }
        currentCluster = []

        X.forEach((dOfX: ILocation) => {
          currentCluster.push(dOfX.rid)
          clusterMap[dOfX.rid] = currentClusterLabel
          Q.push(dOfX)
        })

        while (Q.length > 0) {
          const currentObj: ILocation = <ILocation>Q.pop()
          const Y: ILocation[] = Retrieve_Neighbors(currentObj, D, EpsT, EpsS, embedPointDict)

          if (Y.length >= MinPts) {
            Y.forEach((dOfY: ILocation) => {
              if ((clusterMap[dOfY.rid] !== -1) && !(dOfY.rid in clusterMap)) {

                const avgXY: [number, number] = Cluster_Avg(currentCluster, embedPointDict)
                const valueXY: [number, number] = [embedPointDict[dOfY.rid].left, embedPointDict[dOfY.rid].top] // o.Value
                const xd: number = avgXY[0] - valueXY[0]
                const yd: number = avgXY[1] - valueXY[1]

                if (Math.sqrt(xd*xd + yd*yd) <= deltaE) {
                  currentCluster.push(dOfY.rid)
                  clusterMap[dOfY.rid] = currentClusterLabel
                  Q.push(dOfY)
                }
              }
            })
          }
        }

      }
    }
  })

  return { clusters, noise }
}

function Retrieve_Neighbors(d: ILocation, D: ILocation[], EpsT: number, EpsS: number, embedPointDict: {[rid: number]: ProjPoint}): ILocation[] {
  const ret: ILocation[] = []

  const t: [number, number] = [embedPointDict[d.rid].left, embedPointDict[d.rid].top]

  D.forEach((dOfD: ILocation) => {
    if (dOfD.rid !== d.rid) {
      const sd: number = Get_Geo_Distance(d.lat, d.lng, dOfD.lat, dOfD.lng) // km

      const t2: [number, number] = [embedPointDict[dOfD.rid].left, embedPointDict[dOfD.rid].top]
      const td: number = Math.sqrt((t[0] - t2[0]) * (t[0] - t2[0]) + (t[1] - t2[1]) * (t[1] - t2[1]))

      if (sd <= EpsS && td <= EpsT) {
        ret.push(dOfD)
      }
    }
  })
  return ret
}

function Cluster_Avg(currentCluster: number[], embedPointDict: {[rid: number]: ProjPoint}): [number, number] {
  const ret: [number, number] = [0, 0]
  currentCluster.forEach((rid: number) => {
    const embedPoint = embedPointDict[rid]
    ret[0] += embedPoint.left
    ret[1] += embedPoint.top
  })
  return [ret[0]/currentCluster.length, ret[1]/currentCluster.length]
}

function Get_Geo_Distance(lat1: number, lng1: number, lat2: number, lng2: number) {
	if ((lat1 == lat2) && (lng1 == lng2)) {
		return 0;
	}
	else {
		const radlat1 = Math.PI * lat1/180;
		const radlat2 = Math.PI * lat2/180;
		const theta = lng1-lng2;
		const radtheta = Math.PI * theta/180;
		let dist = Math.sin(radlat1) * Math.sin(radlat2) + Math.cos(radlat1) * Math.cos(radlat2) * Math.cos(radtheta);
		if (dist > 1) {
			dist = 1;
		}
		dist = Math.acos(dist);
		dist = dist * 180/Math.PI;
		dist = dist * 60 * 1.1515;
		dist = dist * 1.609344;
		return dist;
	}
}

// function Get_Value(rid: number, embedPointDict: {[rid: number]: ProjPoint}): [number, number] {
//   return [0, 0]
// }

export { stdbscan, Get_Geo_Distance }
