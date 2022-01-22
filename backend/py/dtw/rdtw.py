import sys
import numpy as np
from sklearn.manifold import TSNE
import matplotlib.pyplot as plt
import json




def rdtw(mat, nr, nc, range_seq_1, range_seq_2):
  mat[0][0]['w'] = mat[0][0]['d']

  # dynamic programing
  for rowIndex, row in enumerate(mat):
    for colIndex, cell in enumerate(row):

      if (colIndex == 0) and (rowIndex == 0):
        continue

      if (rowIndex > 0) and (colIndex > 0):
        v1 = cell['d']*2 + mat[rowIndex-1][colIndex-1]['w']
        v2 = cell['d'] + mat[rowIndex][colIndex-1]['w']
        v3 = cell['d'] + mat[rowIndex-1][colIndex]['w']
        cell['w'] = smallest(v1, v2, v3)

        if (v1 < v2) and (v1 < v3):
          cell['w'] = v1
          cell['p'] = [rowIndex - 1, colIndex - 1]
        elif (v2 < v1) and (v2 < v3):
          cell['w'] = v2
          cell['p'] = [rowIndex, colIndex - 1]
        else:
          cell['w'] = v3
          cell['p'] = [rowIndex - 1, colIndex]

      elif (colIndex == 0) and (rowIndex > 0):
        cell['w'] = cell['d'] + mat[rowIndex-1][colIndex]['w']
        cell['p'] = [rowIndex - 1, colIndex]

      elif (colIndex > 0) and (rowIndex == 0):
        cell['w'] = cell['d'] + mat[rowIndex][colIndex - 1]['w']
        cell['p'] = [rowIndex, colIndex - 1]


  # store the cost
  rowVisitedDict = {}
  for i in range(nr):
    rowVisitedDict[i] = sys.maxsize
  rowVisitedDict[nr - 1] = mat[nr - 1][nc - 1]['d']

  # get shortest path
  pointR = nr - 1
  pointC = nc - 1
  paths = []
  paths.append([pointR, pointC]) # has many point like [2, 3]
  while True:
    currentPoint = paths[-1]
    currentPointR = currentPoint[0]
    currentPointC = currentPoint[1]

    newCell = mat[currentPointR][currentPointC]
    newPointR = newCell['p'][0]
    newPointC = newCell['p'][1]

    if rowVisitedDict[newPointR] > mat[newCell['i'][0]][newCell['i'][1]]['d']:
      rowVisitedDict[newPointR] = mat[newCell['i'][0]][newCell['i'][1]]['d']

    paths.append([newPointR, newPointC])

    # print(newPointR, newPointC, newCell['w'], range_seq_1[newCell['i'][0]], range_seq_2[newCell['i'][1]])

    if newPointR == 0 and newPointC == 0:
      break

  # calculate the cost, referring to the first seq
  distanceBetweenTwoRangeSeqs = 0
  for key, value in rowVisitedDict.items():
    distanceBetweenTwoRangeSeqs += value
  distanceBetweenTwoRangeSeqs = distanceBetweenTwoRangeSeqs/nr

  # print(distanceBetweenTwoRangeSeqs)
  # print(rowVisitedDict)

  return paths, distanceBetweenTwoRangeSeqs








seqs = []
max_seq_len = 0
i = 0
ids = []
for line in open("data/congestionEvents-zj-filter.txt"):
  event_seq_str = line.split('#')[1][:-1].split(',')
  ID = int(line.split('#')[0].split(',')[0])

  # debug
  # if (len(event_seq_str) > 200) or (len(event_seq_str) < 150):
  if len(event_seq_str) < 2:
    # seqs.append([])
    continue

  ids.append(ID)
  if len(event_seq_str) > max_seq_len:
    max_seq_len = len(event_seq_str)
  event_seq_num = list(map(lambda x: int(x), event_seq_str))
  seqs.append(event_seq_num)


dMAT = []

# initialize dMAT
for i, seqi in enumerate(seqs):
  r = []
  for j, seqj in enumerate(seqs):
    r.append(0)
  dMAT.append(r)

for i, seqi in enumerate(seqs):
  print(i, len(seqs))
  for j, seqj in enumerate(seqs):
    if i > j:
      continue

    range_seq_1 = getRangeSeq(seqi)
    range_seq_2 = getRangeSeq(seqj)

    tmp = []
    if len(range_seq_1) > len(range_seq_2):
      tmp = range_seq_1
      range_seq_1 = range_seq_2
      range_seq_2 = tmp

    dMat = getDetailDistanceMatBetweenTwoRangeSeqs(range_seq_1, range_seq_2)
    _, d = rdtw(dMat, len(range_seq_1), len(range_seq_2), range_seq_1, range_seq_2)

    dMAT[i][j] = d
    dMAT[j][i] = d

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

with open('coordinates-dtw.json', 'w+', encoding='UTF-8') as of:
  json.dump(results, of)

plt.scatter(dots[:, 0], dots[:, 1])
plt.show()



# sample_seq_1 = seqs[8]
# sample_seq_2 = seqs[10]
# print(len(sample_seq_1), len(sample_seq_2))

# range_seq_1 = getRangeSeq(sample_seq_1)
# range_seq_2 = getRangeSeq(sample_seq_2)

# # exchange the seq to ensure that the first seq is shorter
# # hence, the row < the col
# tmp = []
# if len(range_seq_1) > len(range_seq_2):
#   tmp = range_seq_1
#   range_seq_1 = range_seq_2
#   range_seq_2 = tmp

# dMat = getDetailDistanceMatBetweenTwoRangeSeqs(range_seq_1, range_seq_2)
# print(len(range_seq_1), len(range_seq_2))
# # print(dMat[0])
# rdtw(dMat, len(range_seq_1), len(range_seq_2), range_seq_1, range_seq_2)

