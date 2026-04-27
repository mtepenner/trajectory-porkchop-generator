# Interplanetary Trajectory "Porkchop Plot" Generator

A high-performance astrodynamics engine and mission planning tool designed to calculate and visualize optimal interplanetary transfer trajectories. This system generates $\Delta v$ contour maps (porkchop plots) and renders 3D flight paths.

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Technologies](#technologies)
- [Installation](#installation)
- [License](#license)

## 🚀 Features
- **High-Speed Lambert Solver**: Implements Izzo's algorithm to parallelize over 10,000+ trajectory solutions for given date ranges.
- **Accurate Ephemeris Data**: Parses NASA JPL SPICE kernels to determine exact planetary positions.
- **Interactive Porkchop Plots**: Renders complex $\Delta v$ heatmap contours to help identify optimal launch and arrival windows.
- **3D Solar System Visualizer**: A rich WebGL interface to visualize the computed transfer orbits between planets.
- **Instantaneous Caching**: Redis-backed caching ensures that heavy matrix computations load instantly on the frontend.

## 🏗️ Architecture
The project utilizes a microservices architecture for massive compute scaling:
1.  **Compute Engine (Go/C++)**: A CPU-optimized gRPC server dedicated to crunching orbital math and state vectors.
2.  **API Gateway (Python/FastAPI)**: Serves as the client-facing API and manages Redis caching for the generated trajectory matrices.
3.  **Mission Planner UI (React)**: The frontend dashboard utilizing D3.js for contour plotting and Three.js for 3D visualization.

## 🛠️ Technologies
- **Astrodynamics**: Go / C++ (Kepler & Lambert Solvers, NASA JPL SPICE)
- **Backend API**: Python, FastAPI, NumPy, Redis
- **Frontend**: React, TypeScript, D3.js (d3-contour), Three.js (@react-three/fiber)
- **Deployment**: Docker Compose, Kubernetes

## 📥 Installation
1. Clone the repository: `git clone https://github.com/mtepenner/trajectory-porkchop-generator.git`
2. Start the local cluster: `docker-compose up`
*(Note: SPICE kernels may need to be downloaded from NASA JPL before running the compute engine).*

## ⚖️ License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
