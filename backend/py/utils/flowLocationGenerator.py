import json
import math

# x_pi = 3.14159265358979324 * 3000.0 / 180.0
# pi = 3.1415926535897932384626  # π
# a = 6378245.0  # 长半轴
# ee = 0.00669342162296594323  # 偏心率平方

# def _transformlat(lng, lat):
#   ret = -100.0 + 2.0 * lng + 3.0 * lat + 0.2 * lat * lat + \
#         0.1 * lng * lat + 0.2 * math.sqrt(math.fabs(lng))
#   ret += (20.0 * math.sin(6.0 * lng * pi) + 20.0 *
#           math.sin(2.0 * lng * pi)) * 2.0 / 3.0
#   ret += (20.0 * math.sin(lat * pi) + 40.0 *
#           math.sin(lat / 3.0 * pi)) * 2.0 / 3.0
#   ret += (160.0 * math.sin(lat / 12.0 * pi) + 320 *
#           math.sin(lat * pi / 30.0)) * 2.0 / 3.0
#   return ret

# def _transformlng(lng, lat):
#   ret = 300.0 + lng + 2.0 * lat + 0.1 * lng * lng + \
#         0.1 * lng * lat + 0.1 * math.sqrt(math.fabs(lng))
#   ret += (20.0 * math.sin(6.0 * lng * pi) + 20.0 *
#           math.sin(2.0 * lng * pi)) * 2.0 / 3.0
#   ret += (20.0 * math.sin(lng * pi) + 40.0 *
#           math.sin(lng / 3.0 * pi)) * 2.0 / 3.0
#   ret += (150.0 * math.sin(lng / 12.0 * pi) + 300.0 *
#           math.sin(lng / 30.0 * pi)) * 2.0 / 3.0
#   return ret

# def gcj02_to_wgs84(lat, lng):
#   """
#   GCJ02(火星坐标系)转GPS84
#   :param lng:火星坐标系的经度
#   :param lat:火星坐标系纬度
#   :return:
#   """
#   # if out_of_china(lng, lat):
#   #   return [lng, lat]
#   dlat = _transformlat(lng - 105.0, lat - 35.0)
#   dlng = _transformlng(lng - 105.0, lat - 35.0)
#   radlat = lat / 180.0 * pi
#   magic = math.sin(radlat)
#   magic = 1 - ee * magic * magic
#   sqrtmagic = math.sqrt(magic)
#   dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * pi)
#   dlng = (dlng * 180.0) / (a / sqrtmagic * math.cos(radlat) * pi)
#   mglat = lat + dlat
#   mglng = lng + dlng
#   return [lat * 2 - mglat, lng * 2 - mglng]


topLeftLatLng = [39.99591903798722, 116.25459597]
bottomRightLatLng = [39.81879326, 116.489552140]

# topLeftLatLng = gcj02_to_wgs84(39.995060447, 116.25459597)
# bottomRightLatLng = gcj02_to_wgs84(39.81879326, 116.489552140)

# topLeftLatLng = gcj02_to_wgs84(39.996874, 116.260289)
# bottomRightLatLng = gcj02_to_wgs84(39.819347, 116.495613)
n = 32
nGrid = 32 * 32

latInterval = (topLeftLatLng[0] - bottomRightLatLng[0]) / 32
lngInterval = (topLeftLatLng[1] - bottomRightLatLng[1]) / 32

topLeftGridCenterLatLng = [
  topLeftLatLng[0] - latInterval/2,
  topLeftLatLng[1] + lngInterval/2
]

locations = []

rid = 0
for i in range(n):
  for j in range(n):
    lat = topLeftGridCenterLatLng[0] - i * latInterval
    lng = topLeftGridCenterLatLng[1] - j * lngInterval
    locations.append({
      "lat": lat,
      "lng": lng,
      "rid": rid
    })
    rid += 1

with open('output/flowData/locations.json', 'w+', encoding='UTF-8') as of:
  json.dump(locations, of)

