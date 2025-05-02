# backend/database/clustering_cuml.py
import os
import json
import logging
from typing import List, Dict, Any
import numpy as np
import pandas as pd

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

try:
    import cudf
    from cuml.feature_extraction.text import TfidfVectorizer
    from cuml.cluster import KMeans

    CUML_AVAILABLE = True
    logger.info("RAPIDS cuML successfully imported")
except ImportError:
    logger.warning("RAPIDS cuML not available, will use CPU fallback")
    CUML_AVAILABLE = False

try:
    from openai import OpenAI

    OPENAI_AVAILABLE = True
except ImportError:
    logger.warning("OpenAI package not available, will use basic cluster descriptions")
    OPENAI_AVAILABLE = False


class ClusterEntry:
    def __init__(self, cluster_description: str, text_entries: List[str]):
        self.cluster_description = cluster_description
        self.text_entries = text_entries

    def to_dict(self) -> Dict[str, Any]:
        return {
            "cluster_description": self.cluster_description,
            "text_entries": self.text_entries,
        }


class TicketResponse:
    def __init__(self, clusters: List[ClusterEntry]):
        self.clusters = clusters

    def to_dict(self) -> Dict[str, Any]:
        return {"clusters": [cluster.to_dict() for cluster in self.clusters]}

    def to_json(self) -> str:
        return json.dumps(self.to_dict())


def generate_cluster_descriptions(clusters: List[List[str]]) -> List[ClusterEntry]:
    """
    Generate descriptions for clusters using OpenAI API or a simple fallback.

    Args:
        clusters: List of clusters, where each cluster is a list of text entries

    Returns:
        List of ClusterEntry objects with descriptions and text entries
    """
    if OPENAI_AVAILABLE:
        try:
            api_key = os.getenv("OPENAI_API_KEY")
            if not api_key:
                logger.error("OPENAI_API_KEY not found in environment variables")
                return [
                    ClusterEntry(f"Cluster {i+1}", cluster)
                    for i, cluster in enumerate(clusters)
                ]

            client = OpenAI(api_key=api_key)
            prompt_engineering = "You are a helpful assistant. Can you provide a short summary of the main features of these products? Write a very short (5 words maximum) title that summarizes all of them collectively."

            cluster_entries = []
            for i, cluster in enumerate(clusters):
                if not cluster:
                    continue

                content = f"Classify the following cluster:\n{json.dumps(cluster)}"

                response = client.chat.completions.create(
                    model="gpt-4o",
                    messages=[
                        {"role": "system", "content": prompt_engineering},
                        {"role": "user", "content": content},
                    ],
                    max_tokens=100,
                )

                description = response.choices[0].message.content.strip()
                cluster_entries.append(ClusterEntry(description, cluster))

            return cluster_entries
        except Exception as e:
            logger.error(f"Error generating cluster descriptions with OpenAI: {str(e)}")
            return [
                ClusterEntry(f"Cluster {i+1}", cluster)
                for i, cluster in enumerate(clusters)
            ]
    else:
        return [
            ClusterEntry(f"Cluster {i+1}", cluster)
            for i, cluster in enumerate(clusters)
        ]


def gpu_tfidf_kmeans_clustering(
    documents: List[str], num_clusters: int = 3
) -> TicketResponse:
    """
    Perform TF-IDF vectorization and K-means clustering using RAPIDS cuML.

    Args:
        documents: List of document texts to cluster
        num_clusters: Number of clusters to create

    Returns:
        TicketResponse object containing the clusters and their descriptions
    """
    if not documents:
        logger.warning("No documents provided for clustering")
        return TicketResponse([])

    try:
        documents_series = cudf.Series(documents)

        logger.info("Performing TF-IDF vectorization on GPU")
        vectorizer = TfidfVectorizer(stop_words="english")
        tfidf_matrix = vectorizer.fit_transform(documents_series)

        logger.info(f"Performing K-means clustering with {num_clusters} clusters")
        kmeans = KMeans(n_clusters=num_clusters, random_state=0)
        cluster_labels = kmeans.fit_predict(tfidf_matrix)

        cluster_labels_cpu = cluster_labels.to_numpy()

        clustered_texts = [[] for _ in range(num_clusters)]
        for idx, label in enumerate(cluster_labels_cpu):
            clustered_texts[label].append(documents[idx])

        cluster_entries = generate_cluster_descriptions(clustered_texts)

        return TicketResponse(cluster_entries)

    except Exception as e:
        logger.error(f"Error in GPU-accelerated clustering: {str(e)}")
        return cpu_tfidf_kmeans_clustering(documents, num_clusters)


def cpu_tfidf_kmeans_clustering(
    documents: List[str], num_clusters: int = 3
) -> TicketResponse:
    """
    CPU-based implementation of TF-IDF and K-means clustering using scikit-learn.

    This function is used when GPU processing is not available or fails.

    Args:
        documents: List of document texts to cluster
        num_clusters: Number of clusters to create

    Returns:
        TicketResponse object containing the clusters and their descriptions
    """
    try:
        from sklearn.feature_extraction.text import (
            TfidfVectorizer as SklearnTfidfVectorizer,
        )
        from sklearn.cluster import KMeans as SklearnKMeans

        logger.info("Using CPU-based TF-IDF and K-means clustering")

        vectorizer = SklearnTfidfVectorizer(stop_words="english")
        tfidf_matrix = vectorizer.fit_transform(documents)

        kmeans = SklearnKMeans(n_clusters=num_clusters, random_state=0)
        cluster_labels = kmeans.fit_predict(tfidf_matrix)

        clustered_texts = [[] for _ in range(num_clusters)]
        for idx, label in enumerate(cluster_labels):
            clustered_texts[label].append(documents[idx])

        cluster_entries = generate_cluster_descriptions(clustered_texts)

        return TicketResponse(cluster_entries)
    except Exception as e:
        logger.error(f"Error in CPU-based clustering: {str(e)}")
        return TicketResponse([ClusterEntry("Error in clustering", documents)])


def tfidf_kmeans_clustering_api(documents: List[str]) -> str:
    """
    API function to be called from Go code.

    Args:
        documents: List of document texts to cluster

    Returns:
        JSON string representation of the clustering result
    """
    if CUML_AVAILABLE:
        response = gpu_tfidf_kmeans_clustering(documents)
    else:
        response = cpu_tfidf_kmeans_clustering(documents)

    return response.to_json()


if __name__ == "__main__":
    test_documents = [
        "Network issue on east datacenter",
        "Cannot connect to database server",
        "Email service is down",
        "Server maintenance required",
        "Website loading slowly",
        "Database backup failed",
        "Network connectivity problems",
        "Email delivery delays",
        "Server CPU at 100%",
        "Website shows 404 errors",
    ]

    if CUML_AVAILABLE:
        print("Testing GPU implementation:")
        result = gpu_tfidf_kmeans_clustering(test_documents)
    else:
        print("Testing CPU implementation:")
        result = cpu_tfidf_kmeans_clustering(test_documents)

    print(json.dumps(result.to_dict(), indent=2))
