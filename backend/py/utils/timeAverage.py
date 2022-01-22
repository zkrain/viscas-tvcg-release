import json
import numpy as np
import datetime


def runAir(eventFilepath, saveFilePath, rowNum=12, timeDelta=1, startDateTime=datetime.datetime(2018, 1, 1, 0, 0, 0)):
    eventFile = open(eventFilepath, 'r')
    lines = eventFile.readlines()
    res = {}
    for line in lines:
        eventNums = np.zeros(rowNum)
        durations = np.zeros(rowNum)
        infectNum = np.zeros(rowNum)
        head, eventStr = line.split('#')
        rid = int(head.split(',')[0])
        eventsStr = eventStr.split(',')
        events = [int(s) for s in eventsStr]
        for event in events:
            h = datetime.timedelta(hours=event * timeDelta) + startDateTime
            ind = h.month
            eventNums[ind - 1] += 1
        eventNums[11] = int(eventNums[11] * 31 / 23)  # 因为12月的数据只有23天，为了平衡，就乘31/23
        _, infectEvents = eventsToInfectEvent(events)
        print(infectEvents)
        for infectEvent in infectEvents:
            start = infectEvent[0]
            h = datetime.timedelta(hours=start * timeDelta) + startDateTime
            ind = h.month
            durations[ind - 1] = durations[ind - 1] + (infectEvent[1] - infectEvent[0] + 1)
            infectNum[ind - 1] += 1
        aveDuration = []
        for s, n in zip(durations, infectNum):
            if n > 0:
                aveDuration += [s / n]
            else:
                aveDuration += [0]
        res[rid] = [[n, a] for n, a in zip(eventNums, aveDuration)]
    with open(saveFilePath, 'w') as file:
        json.dump(res, file)


def runCongestion(eventFilepath, saveFilePath, rowNum=24, timeDeltaMinutes=10,
                  startDateTime=datetime.datetime(2018, 3, 1, 0, 0, 0)):
    eventFile = open(eventFilepath, 'r')
    lines = eventFile.readlines()
    res = {}
    for line in lines:
        eventNums = np.zeros(rowNum)
        durations = np.zeros(rowNum)
        infectNum = np.zeros(rowNum)
        head, eventStr = line.split('#')
        rid = int(head.split(',')[0])
        eventsStr = eventStr.split(',')
        events = [int(s) for s in eventsStr]
        for event in events:
            h = datetime.timedelta(minutes=event * timeDeltaMinutes) + startDateTime
            ind = h.hour
            eventNums[ind - 1] += 1
        _, infectEvents = eventsToInfectEvent(events)
        for infectEvent in infectEvents:
            start = infectEvent[0]
            h = datetime.timedelta(minutes=start * timeDeltaMinutes) + startDateTime
            ind = h.hour
            durations[ind - 1] = durations[ind - 1] + (infectEvent[1] - infectEvent[0] + 1)
            infectNum[ind - 1] += 1
        aveDuration = []
        for s, n in zip(durations, infectNum):
            if n > 0:
                aveDuration += [s / n]
            else:
                aveDuration += [0]
        res[rid] = [[n, a] for n, a in zip(eventNums, aveDuration)]
    print(res)
    with open(saveFilePath, 'w') as file:
        json.dump(res, file)


def eventsToInfectEvent(events, deltaTime=1):
    infectEvents = []
    infectEventNum = 0
    if len(events) == 0:
        return infectEventNum, infectEvents
    startTime = events[0]
    endTime = events[0]
    for i in range(1, len(events)):
        if events[i] - events[i - 1] > deltaTime + 1:
            infectEvents += [[startTime, endTime]]
            infectEventNum += 1
            startTime = events[i]
        endTime = events[i]
        if i == len(events) - 1:
            infectEvents += [[startTime, endTime]]
            infectEventNum += 1
    return infectEventNum, infectEvents


def test():
    a = [1, 2, 3, 5, 8, 9, 11, 16]
    print(eventsToInfectEvent(a))


if __name__ == '__main__':
    # test()
    filePath = '../../output/flowData/flowEvents.txt'
    savePath = '../../output/flowData/flowAverages.json'
    runCongestion(filePath, savePath, timeDeltaMinutes=30, startDateTime=datetime.datetime(2015, 11, 1, 0, 0, 0))
    # filePath = '../../output/airData/airEvents.txt'
    # savePath = '../../output/airData/airEventsAverages.json'
    # runAir(filePath, savePath)
