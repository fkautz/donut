package client

import (
	"sort"
	"strings"
)

func filterPrefix(objects []string, prefix string) []string {
	var results []string
	for _, object := range objects {
		if strings.HasPrefix(object, prefix) {
			results = append(results, object)
		}
	}
	return results
}

func removePrefix(objects []string, prefix string) []string {
	var results []string
	for _, object := range objects {
		results = append(results, strings.TrimPrefix(object, prefix))
	}
	return results
}

func filterDelimited(objects []string, delim string) []string {
	var results []string
	for _, object := range objects {
		if !strings.Contains(object, delim) {
			results = append(results, object)
		}
	}
	return results
}
func filterNotDelimited(objects []string, delim string) []string {
	var results []string
	for _, object := range objects {
		if strings.Contains(object, delim) {
			results = append(results, object)
		}
	}
	return results
}

func extractDir(objects []string, delim string) []string {
	var results []string
	for _, object := range objects {
		parts := strings.Split(object, delim)
		results = append(results, parts[0]+"/")
	}
	return results
}

func uniqueObjects(objects []string) []string {
	objectMap := make(map[string]string)
	for _, v := range objects {
		objectMap[v] = v
	}
	var results []string
	for k := range objectMap {
		results = append(results, k)
	}
	sort.Strings(results)
	return results
}
