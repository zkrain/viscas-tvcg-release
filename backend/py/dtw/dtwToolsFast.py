import sys

def getDistanceBetweenTwoRanges(r1, r2):
  # return (abs(r1[0] - r2[0]) + abs(r1[1] - r2[1]))/(r1[2] + r2[2])
  return (abs(r1[2] - r2[2]) + abs(r1[0] - r2[0]))/(r1[2] + r2[2])

def getRangeSeq(seq):
  rangeSeq = []
  currentRangeEvent = [-1, -1, 0] # start, end, length
  for index, t in enumerate(seq):
    if currentRangeEvent[2] == 0:
      # precess a new range
      currentRangeEvent[2] += 1
      currentRangeEvent[1] = t + 1
      currentRangeEvent[0] = t
    else:
      # precess an old range
      if currentRangeEvent[1] != t:
        # pop
        rangeSeq.append(currentRangeEvent)
        currentRangeEvent = [-1, -1, 0]
      else:
        currentRangeEvent[2] += 1
        currentRangeEvent[1] = t + 1

  return rangeSeq

def refineDistance(path, range_seq_1, range_seq_2):
  nr = len(range_seq_1)
  rowVisitedDict = {}
  for i in range(nr):
    rowVisitedDict[i] = sys.maxsize

  for p in path:
    range_of_1 = range_seq_1[p[0]]
    range_of_2 = range_seq_2[p[1]]
    d_of_ranges = getDistanceBetweenTwoRanges(range_of_1, range_of_2)
    if rowVisitedDict[p[0]] > d_of_ranges:
      rowVisitedDict[p[0]] = d_of_ranges

  refinedDistance = 0
  for key, value in rowVisitedDict.items():
    refinedDistance += value
  refinedDistance = refinedDistance/nr

  return refinedDistance
