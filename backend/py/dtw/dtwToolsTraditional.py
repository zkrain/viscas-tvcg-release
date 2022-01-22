# traditional dtw tools
def getDetailDistanceMatBetweenTwoRangeSeqs(rSeq1, rSeq2):
  mat = []
  for index1, range1 in enumerate(rSeq1):
    row = []
    for index2, range2 in enumerate(rSeq2):
      d = getDistanceBetweenTwoRanges(range1, range2)
      cell = {
        'd': d, # static
        'p': [0, 0], # the parent location
        'w': 0, # can be updated
        'i': [index1, index2]
      }
      row.append(cell)
    mat.append(row)
  return mat


def getDistanceBetweenTwoRanges(r1, r2):
  return abs(r1['s'] - r2['s']) + abs(r1['l'] - r2['l'])

def getRangeSeq(seq):
  rangeSeq = []
  currentRangeEvent = {
    's': -1, # start
    'e': -1, # end
    'l': 0
  }
  for index, t in enumerate(seq):
    if currentRangeEvent['l'] == 0:
      # precess a new range
      currentRangeEvent['l'] += 1
      currentRangeEvent['e'] = t + 1
      currentRangeEvent['s'] = t
    else:
      # precess an old range
      if currentRangeEvent['e'] != t:
        # pop
        rangeSeq.append(currentRangeEvent)
        currentRangeEvent = {
          's': -1, # start
          'e': -1, # end
          'l': 0
        }
      else:
        currentRangeEvent['l'] += 1
        currentRangeEvent['e'] = t + 1

  return rangeSeq

def smallest(num1, num2, num3):
  if (num1 < num2) and (num1 < num3):
    return num1
  elif (num2 < num1) and (num2 < num3):
    return num2
  else:
    return num3

