from sklearn.cluster import DBSCAN
import numpy as np
import json
import matplotlib.pyplot as plt

coordinates = []
with open("output/coordinates-dtw.json",'r') as load_f:
  coordinates = json.load(load_f)

coordinatesArrFomrat = []
for coordinate in coordinates:
  coordinatesArrFomrat.append([coordinate['x'], coordinate['y']])

X = np.array(coordinatesArrFomrat)
clustering = DBSCAN(eps=0.175, min_samples=12, algorithm="kd_tree", p=2.5).fit(X)
# 0.21 15

# for label in clustering.labels_:
#   print(label)
group = clustering.labels_
# print(type(clustering.labels_))


plt.figure(figsize=(8, 6))
plt.scatter(X[:, 0], X[:, 1], c=clustering.labels_.astype(float))

# fig, ax = plt.subplots()
# for g in np.unique(group):
#   ix = np.where(group == g)
#   ax.scatter(scatter_x[ix], scatter_y[ix], c = cdict[g], label = g, s = 100)
# ax.legend()
plt.show()



