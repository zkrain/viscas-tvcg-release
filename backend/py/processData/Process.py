import h5py
import numpy as np
import datetime
from numpy.core.multiarray import ndarray


def exploreData():
    path = 'D:/codes/vis2020/Data/TaxiBJ/BJ16_M32x32_T30_InOut.h5'
    f = h5py.File(path)
    date: ndarray = np.array(f['date'])
    for d in date:
        print(d)
    print(date)
    print(date.shape[0])


def run(inputPath: str, outputPath: str):
    f = h5py.File(inputPath)
    data: ndarray = np.array(f['data'])
    date: ndarray = np.array(f['date'])
    for d in date:
        print(d)
    print(date)
    print(data.shape)
    flow_dict: dict = {}
    ts, _, m, n = data.shape
    for t in range(ts):
        for i in range(m):
            for j in range(n):
                ind: int = 32 * i + j
                if ind not in flow_dict.keys():
                    flow_dict[ind] = [data[t][0][i][j] - data[t][1][i][j]]
                else:
                    flow_dict[ind] += [data[t][0][i][j] - data[t][1][i][j]]
    start = byte2date(date[0])
    file = open(outputPath, 'w')
    for locationId in flow_dict.keys():
        flow = flow_dict[locationId]
        eventsStr = ''
        length = 0
        for fi in range(len(flow)):
            if flow[fi] > 2:
                delta = (byte2date(date[fi]) - start)
                stamp = int(delta.total_seconds() / (30 * 60))
                eventsStr = eventsStr + ',' + str(stamp)
                length += 1
        if length > 50:
            rowStr = str(locationId) + ',' + str(length) + ',999#' + eventsStr[1:] + '\n'
            file.write(rowStr)
    file.close()

def byte2date(b):
    dateStr = str(b, encoding='utf-8')
    return datetime.datetime.strptime(dateStr[:-2], '%Y%m%d') + datetime.timedelta(hours=(float(dateStr[-2:]) - 1)/2)



if __name__ == '__main__':
    inputPath = 'D:/codes/vis2020/Data/TaxiBJ/BJ16_M32x32_T30_InOut.h5'
    outputPath = 'D:/codes/vis2020/Data/TaxiBJ/output/BJ16InflowEvents.txt'
    run(inputPath, outputPath)
