package languages

import (
	"github.com/hashicorp/go-version"
	"sort"
)

func SortVersions(versionsRaw []string) ([]string, error) {
	versions := make([]*version.Version, len(versionsRaw))
	for i, raw := range versionsRaw {
		v, _ := version.NewVersion(raw)

		versions[i] = v
	}

	// After this, the versions are sorted in descending order properly sorted
	sort.Sort(sort.Reverse(version.Collection(versions)))

	var versionStrings []string
	for _, v := range versions {
		versionStrings = append(versionStrings, v.String())
	}
	return versionStrings, nil
}

func removeElements(originalElements []string, elementsToRemove []string) []string {
	// Iterate over the slice to be modified
	for i := 0; i < len(originalElements); i++ {
		// Iterate over the elements to be removed
		for _, elem := range elementsToRemove {
			if originalElements[i] == elem {
				// Remove the matching element from the slice
				originalElements = append(originalElements[:i], originalElements[i+1:]...)
				i--   // Decrement the index to account for the removed element
				break // Break out of the inner loop to avoid removing the same element twice
			}
		}
	}
	return originalElements
}
