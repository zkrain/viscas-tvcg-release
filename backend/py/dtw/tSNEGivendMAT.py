import json
import numpy as np
from sklearn.manifold import TSNE
import matplotlib.pyplot as plt

dMAT = {}
with open('py/output-flow/dMAT-dtw.json','r') as load_f:
  dMAT = json.load(load_f)

ids = []
with open("py/output-flow/ids-dtw.json",'r') as load_f:
  ids = json.load(load_f)

print(type(dMAT))
dMATnp = np.array(dMAT)
dots = TSNE(n_components=2, metric="precomputed", random_state=232, n_iter=2000).fit_transform(dMATnp)
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

with open('py/output-flow/coordinates-2-dtw.json', 'w+', encoding='UTF-8') as of:
  json.dump(results, of)


plt.scatter(dots[:, 0], dots[:, 1])
plt.show()