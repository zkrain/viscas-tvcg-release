# Visual Cascade Analytics of Large-scale Spatiotemporal Data

[Paper Link](https://zjuidg.org/source/projects/VisCas/VisCas.pdf)

Abstract: Many spatiotemporal events can be viewed as contagions. These events implicitly propagate across space and time by following cascading patterns, expanding their influence, and generating event cascades that involve multiple locations. Analyzing such cascading processes presents valuable implications in various urban applications, such as traffic planning and pollution diagnostics. Motivated by the limited capability of the existing approaches in mining and interpreting cascading patterns, we propose a visual analytics system called VisCas. VisCas combines an inference model with interactive visualizations and empowers analysts to infer and interpret the latent cascading patterns in the spatiotemporal context. To develop VisCas, we address three major challenges, 1) generalized pattern inference, 2) implicit influence visualization, and 3) multifaceted cascade analysis. For the first challenge, we adapt the state-of-the-art cascading network inference technique to general urban scenarios, where cascading patterns can be reliably inferred from large-scale spatiotemporal data. For the second and third challenges, we assemble a set of effective visualizations to support location navigation, influence inspection, and cascading exploration, and facilitate the in-depth cascade analysis. We design a novel influence view based on a three-fold optimization strategy for analyzing the implicit influences of the inferred patterns. We demonstrate the capability and effectiveness of VisCas with two case studies conducted on real-world traffic congestion and air pollution datasets with domain experts.


## Installation

### Prerequisites
* [Docker](https://docs.docker.com/get-docker/) 20.10
* [Docker Compose](https://docs.docker.com/compose/cli-command/#installing-compose-v2) v2
* Google Chrome

### Build & Run
``sudo`` can be omitted if the current user has sufficient privilege to execute docker commands:

```shell script
$ sudo docker compose up
```

VisCas can now be accessed via http://localhost:8080 with Chrome.

## Data description

The datasets directly used by the system are placed under /backend/output/congestion2Data (traffic congestion dataset) and /backend/output/airData (air pollution event dataset).

The traffic congestion dataset is extracted from taxi trajectories that provided by local authorities.
After extraction, the dataset comprises the congestion event series of every road segments.
Due to the sensitivity of the data, onl the **samples** related to the case are released in this repository.

The air pollution event dataset can be collected and derived from air quality websites, for example, [New York Air Pollution](https://aqicn.org/city/usa/newyork/).


## Replicability

We provide a [DEMO Video](https://www.youtube.com/watch?v=IVSf0BNRC_c&t=3s) to guide the reproduction of the cases in [our paper](https://zjuidg.org/source/projects/VisCas/VisCas.pdf).
The video content is divided into three parts.
The first part introduces the system.
The second and third parts introduce two cases, respectively, in a step-by-step way.
After the frontend and backend run, users can follow the steps in the video to reproduce the two cases.

## Note
We use leaflet.js and mapbox for map services. which is for academic purposes **only** and does not indicate any political standpoints.

## Citation
If you use this code for your research, please consider citing:
```
@article{deng2021viscas,
  title={Visual Cascade Analytics of Large-scale Spatiotemporal Data},
  author={Deng, Zikun and Weng, Di and Liang, Yuxuan and Bao, Jie and Zheng, Yu and Schreck, Tobias and Xu, Mingliang and Wu, Yingcai},
  journal={IEEE Transactions on Visualization and Computer Graphics},
  year={2022},
}
```
