from gensim.models import Word2Vec
from sklearn.decomposition import PCA
from sklearn import manifold
from matplotlib import pyplot as plt
import json
import numpy as np
sentences = []

# f = open("output/document.txt")
# line = f.readline()
# while line:
#     sentence = line.replace("\n","").split(" ")
#     sentences.append(sentence)
#     line = f.readline()
# f.close()

with open('output/document.txt', 'r', encoding = 'utf-8') as f:
    for line in f:
        sentence = line.replace("\n","").split(" ")
        sentences.append(sentence)

print("document loaded")

model = Word2Vec(sentences, min_count=5, sg=1, window=3, size=50)

print("model trained")

X = model[model.wv.vocab]
np.savetxt("output/vector.csv", X, delimiter=",")
# pca = PCA(n_components=2)
ebd = manifold.TSNE(n_components=2, n_iter=10000)
result = ebd.fit_transform(X)

plt.scatter(result[:, 0], result[:, 1])
words = list(model.wv.vocab)

'''
Labels

# for i, word in enumerate(words):
#     plt.annotate(word, xy=(result[i, 0], result[i, 1]))

# plt.show()
'''

OutJSON = []

print(len(words))
for i, word in enumerate(words):
    OutJSON.append({'rid': int(word), 'x': float(result[i, 0]), 'y': float(result[i, 1])})

with open('output/coordinates.json', 'w+', encoding='UTF-8') as of:
    json.dump(OutJSON, of)

plt.show()
print(ebd.n_iter_)
