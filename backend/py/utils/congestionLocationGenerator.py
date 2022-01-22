import json

locations = []

with open('output/roadsShape-zj-filter.txt') as f:
  while True:
    line = f.readline()
    if not line:
      break
    rid = line.split('#')[0].split(',')[0]
    shape = line.split('#')[1][:-2].split(';')
    median = shape[round(len(shape)/2)]
    lat = float(median.split(',')[0])
    lng = float(median.split(',')[1])
    locations.append({
      'lat': lat,
      'lng': lng,
      'rid': int(rid)
    })

print(len(locations))

with open('output/congestionData/locations.json', 'w+', encoding='UTF-8') as of:
  json.dump(locations, of)