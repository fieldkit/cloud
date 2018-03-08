package main

import (
	"github.com/hongshibao/go-kdtree"
)

const (
	NOISE     = false
	CLUSTERED = true
)

type Clusterable interface {
	kdtree.Point
	GetID() string
}

type Cluster []Clusterable

func Clusterize(objects []Clusterable, tree *kdtree.KDTree, minPts int, eps float64) []Cluster {
	clusters := make([]Cluster, 0)
	visited := map[string]bool{}
	for _, point := range objects {
		neighbours := findNeighbours(point, objects, tree, eps)
		if len(neighbours)+1 >= minPts {
			visited[point.GetID()] = CLUSTERED
			cluster := make(Cluster, 1)
			cluster[0] = point
			cluster = expandCluster(cluster, neighbours, visited, tree, minPts, eps)

			if len(cluster) >= minPts {
				clusters = append(clusters, cluster)
			}
		} else {
			visited[point.GetID()] = NOISE
		}
	}
	return clusters
}

//Finds the neighbours from given array
//depends on Eps variable, which determines
//the distance limit from the point
func findNeighbours(point Clusterable, points []Clusterable, tree *kdtree.KDTree, eps float64) []Clusterable {
	neighbours := make([]Clusterable, 0)
	for _, potNeigb := range points {
		if point.GetID() != potNeigb.GetID() && potNeigb.Distance(point) <= eps {
			neighbours = append(neighbours, potNeigb)
		}
	}
	return neighbours
}

//Try to expand existing clutser
func expandCluster(cluster Cluster, neighbours []Clusterable, visited map[string]bool, tree *kdtree.KDTree, minPts int, eps float64) Cluster {
	seed := make([]Clusterable, len(neighbours))
	copy(seed, neighbours)
	for _, point := range seed {
		pointState, isVisited := visited[point.GetID()]
		if !isVisited {
			currentNeighbours := findNeighbours(point, seed, tree, eps)
			if len(currentNeighbours)+1 >= minPts {
				visited[point.GetID()] = CLUSTERED
				cluster = merge(cluster, currentNeighbours)
			}
		}

		if isVisited && pointState == NOISE {
			visited[point.GetID()] = CLUSTERED
			cluster = append(cluster, point)
		}
	}

	return cluster
}

func merge(one []Clusterable, two []Clusterable) []Clusterable {
	mergeMap := make(map[string]Clusterable)
	putAll(mergeMap, one)
	putAll(mergeMap, two)
	merged := make([]Clusterable, 0)
	for _, val := range mergeMap {
		merged = append(merged, val)
	}

	return merged
}

//Function to add all values from list to map
//map keys is then the unique collecton from list
func putAll(m map[string]Clusterable, list []Clusterable) {
	for _, val := range list {
		m[val.GetID()] = val
	}
}
