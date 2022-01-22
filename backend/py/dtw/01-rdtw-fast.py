from dtwToolsFast import getRangeSeq, getDistanceBetweenTwoRanges, refineDistance
from fastdtw import fastdtw
import numpy as np
import sys
from sklearn.manifold import TSNE
import matplotlib.pyplot as plt
from multiprocessing import Pool
from functools import partial
import json

from dtwParallelWorker import second_level_parallel

# # congestion
# hackOutlierIDs = {}
# inputFileName = "output/congestionData/congestionEvents.txt"
# outputFileNameIds = "py/output-congestion/ids-dtw.json"
# outputFileNameDMAT = 'py/output-congestion/dMAT-dtw.json'
# outputFileNameDMATRefined = 'py/output-congestion/dMAT-refined-dtw.json'
# outputFileNameCoor = 'py/output-congestion/coordinates-dtw.json'

# # air
hackOutlierIDs = {}
inputFileName = "D:/codes/vis2020/Data/TaxiBJ/output/BJ16InflowEventsBackup.txt"
outputFileNameIds = "D:/codes/vis2020/Data/TaxiBJ/output/ids-dtw.json"
outputFileNameDMAT = 'D:/codes/vis2020/Data/TaxiBJ/output/dMAT-dtw.json'
outputFileNameDMATRefined = 'D:/codes/vis2020/Data/TaxiBJ/output/dMAT-refined-dtw.json'
outputFileNameCoor = 'D:/codes/vis2020/Data/TaxiBJ/output/coordinates-dtw.json'

# flow
# hackOutlierIDs = {"441": 1, "409": 1, "400": 1, "370": 1, "574": 1, "584": 1, "606": 1, "323": 1}
# inputFileName = "output/flowData/flowEvents.txt"
# outputFileNameIds = "py/output-flow/ids-dtw.json"
# outputFileNameDMAT = 'py/output-flow/dMAT-dtw.json'
# outputFileNameDMATRefined = 'py/output-flow/dMAT-refined-dtw.json'
# outputFileNameCoor = 'py/output-flow/coordinates-dtw.json'

seqs = []
max_seq_len = 0
i = 0
ids = []
for line in open(inputFileName):
  event_seq_str = line.split('#')[1][:-1].split(',')
  ID = int(line.split('#')[0].split(',')[0])

  if str(ID) in hackOutlierIDs:
    continue

  # if len(event_seq_str) < 300 or len(event_seq_str) > 1800 :
  #   continue

  ids.append(ID)
  if len(event_seq_str) > max_seq_len:
    max_seq_len = len(event_seq_str)
  event_seq_num = list(map(lambda x: int(x), event_seq_str))
  seqs.append(getRangeSeq(event_seq_num))


with open(outputFileNameIds, 'w+', encoding='UTF-8') as of:
  json.dump(ids, of)


dMAT = []
dMATRefined = []
# initialize dMAT
for i, range_seq_1 in enumerate(seqs):
  r = []
  r2 = []
  for j, range_seq_2 in enumerate(seqs):
    r.append(0)
    r2.append(0)
  dMAT.append(r)
  dMATRefined.append(r)

# # single processing
for i, range_seq_1 in enumerate(seqs):
  print(i, len(seqs))
  for j, range_seq_2 in enumerate(seqs):
    if i > j:
      continue

    # range_seq_1 = getRangeSeq(seqi)
    # range_seq_2 = getRangeSeq(seqj)

    # exchange the seq to ensure that the first seq is shorter
    # hence, the row < the col
    tmp = []
    if len(range_seq_1) > len(range_seq_2):
      tmp = range_seq_1
      range_seq_1 = range_seq_2
      range_seq_2 = tmp

    # range_seq_1 ： [[1,3], [5,4], ...]
    distance, path = fastdtw(range_seq_1, range_seq_2, dist=getDistanceBetweenTwoRanges)
    # distanceRefined = refineDistance(path, range_seq_1, range_seq_2)

    dMAT[i][j] = distance
    dMAT[j][i] = distance
    # dMATRefined[i][j] = distanceRefined
    # dMATRefined[j][i] = distanceRefined
# # single processing end

with open(outputFileNameDMAT, 'w+', encoding='UTF-8') as of:
  json.dump(dMAT, of)
# with open(outputFileNameDMATRefined, 'w+', encoding='UTF-8') as of2:
  # json.dump(dMATRefined, of2)

dMATnp = np.array(dMAT)
dots = TSNE(n_components=2, metric="precomputed").fit_transform(dMATnp)
print(dots.shape)


results = []
for index, ID in enumerate(ids):
  x = dots[index][0]
  y = dots[index][1]
  results.append({
    'x': float(x),
    'y': float(y),
    'rid': ID
  })

with open(outputFileNameCoor, 'w+', encoding='UTF-8') as of:
  json.dump(results, of)

plt.scatter(dots[:, 0], dots[:, 1])
plt.show()
plt.savefig('D:/codes/vis2020/Data/TaxiBJ/output/2.png')