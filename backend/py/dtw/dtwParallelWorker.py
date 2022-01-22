from dtwToolsFast import getDistanceBetweenTwoRanges, refineDistance
from fastdtw import fastdtw


def second_level_parallel(range_seq_2_tuple, range_seq_1, dMAT, i):
  j, range_seq_2 = range_seq_2_tuple
  # exchange the seq to ensure that the first seq is shorter
  # hence, the row < the col
  tmp = []
  if len(range_seq_1) > len(range_seq_2):
    tmp = range_seq_1
    range_seq_1 = range_seq_2
    range_seq_2 = tmp

  _, path = fastdtw(range_seq_1, range_seq_2, dist=getDistanceBetweenTwoRanges)
  distance = refineDistance(path, range_seq_1, range_seq_2)

  dMAT[i][j] = distance
  dMAT[j][i] = distance