import json
import sys

sys.path.append('../..')

def eventsToInfectEvent(events, deltaTime = 1):
    infectEvents = []
    infectEventNum = 0
    if len(events) == 0:
        return infectEventNum, infectEvents
    startTime = events[0]
    endTime = events[0]
    for i in range(1, len(events)):
        if events[i] - events[i - 1] > deltaTime + 1:
            infectEvents += list(range(startTime, endTime + 1))
            infectEventNum += 1
            startTime = events[i]
        endTime = events[i]
        if i == len(events) - 1:
            infectEvents += list(range(startTime, endTime + 1))
            infectEventNum += 1
    return infectEventNum, infectEvents


def run(eventFilePath, outputPath):
    outputFile = open(outputPath, 'w')
    sum1 = 0
    sum2 = 0
    with open(eventFilePath, 'r') as file:
        lines = file.readlines()
        for line in lines:
            head, eventsStr = line.split('#')
            infos = head.split(',')
            rid = infos[0]
            events = [int(e) for e in eventsStr.split(',')]
            infectEventNum, infectEvents = eventsToInfectEvent(events)
            infectEventsStr = ""
            for i, e in enumerate(infectEvents):
                infectEventsStr += str(e)
                if i < len(infectEvents) - 1:
                    infectEventsStr += ","
                else:
                    infectEventsStr += "\n"
            rowStr = str(rid) + "," + str(len(infectEvents)) + "," + "999" + "," + str(infectEventNum) + "#" + infectEventsStr
            sum1 += len(events)
            sum2 += len(infectEvents)
            outputFile.write(rowStr)
    print(sum2 - sum1, len(lines), (sum2 - sum1) / len(lines))
    outputFile.close()

            

if __name__ == "__main__":
    eventFilePath = 'D:/codes/vis2020/Data/TaxiBJ/output/BJ16InflowEvents.txt'
    p1, p2 = eventFilePath.split('.')
    outputPath = p1 + 'Backup.' + p2
    run(eventFilePath, outputPath)
    # print(eventsToInfectEvent([1,2,3,6,8,9])) # test