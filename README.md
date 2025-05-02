# Mimir

Have someone that knows what you need before you even need to ask.

# Inspiration

The inspiration for Mimir came from the need to streamline the management of knowledge and data, helping teams collaborate more effectively on complex technical projects. We wanted to create a solution that leverages automation and intelligence to simplify the process of accessing relevant information quickly, enabling better decision-making and smoother project execution.

# What it does

Mimir is a comprehensive platform that centralizes technical documentation, project updates, and resources, providing personalized recommendations based on user roles and project needs. It integrates with various tools and platforms, offering real-time insights and automated notifications. With advanced search capabilities, it helps users quickly locate critical information, while also recommending resources based on ongoing projects.

## GPU Acceleration 
[![NVIDIA RAPIDS](https://img.shields.io/badge/NVIDIA_RAPIDS-76B900?style=for-the-badge&logo=nvidia&logoColor=white)](https://rapids.ai/)


Mimír uses RAPIDS cuML to accelerate the clustering of incident reports, providing up to 10-50x performance improvements for large datasets. Key machine learning components include:

- **TF-IDF Vectorization**: Converting text descriptions into numerical vectors
- **K-means Clustering**: Grouping similar incidents together

The integration is achieved through a Go-Python bridge `backends/database/clustering_cuml_bridge.go`, which will later on be swapped for gRPC calls.
# How we built it

Using Go, Next.js, Python, cuML, React, Typescript, WeviateDB.

# Flowchart

![Mimir flowchart](https://github.com/user-attachments/assets/ea442464-81e5-40b9-8242-c3c9d1375739)
