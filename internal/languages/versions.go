package languages

import (
	"errors"
	"github.com/hashicorp/go-version"
	"sort"
)

func SortVersionsDesc(versions []*version.Version) []*version.Version {
	// After this, the versions are sorted in descending order properly sorted
	sort.Sort(sort.Reverse(version.Collection(versions)))

	return versions
}

func removeElements(originalElements []string, elementsToRemove []string) ([]string, error) {
	if originalElements == nil || elementsToRemove == nil {
		return nil, errors.New("slice and elements arguments cannot be nil")
	}
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
	return originalElements, nil
}

func removeNewerVersions(versions []*version.Version, maxVersion string) ([]*version.Version, error) {
	maxV, err := version.NewVersion(maxVersion)
	if err != nil {
		return nil, err
	}
	var result []*version.Version
	for _, v := range versions {
		if v.Segments()[0] == maxV.Segments()[0] && v.Segments()[1] == maxV.Segments()[1] && v.Segments()[2] < maxV.Segments()[2] {
			// Version is newer across Major.Minor versions
			result = append(result, v)
		} else if v.Segments()[0] != maxV.Segments()[0] || v.Segments()[1] != maxV.Segments()[1] {
			// Version is newer only within the same Major.Minor version line
			result = append(result, v)
		}
	}
	return result, nil
}

func removeOlderVersions(versions []*version.Version, minVersion string) ([]*version.Version, error) {
	minV, err := version.NewVersion(minVersion)
	if err != nil {
		return nil, err
	}
	var result []*version.Version
	for _, v := range versions {
		if v.Segments()[0] == minV.Segments()[0] && v.Segments()[1] == minV.Segments()[1] && v.Segments()[2] > minV.Segments()[2] {
			// Version is newer across Major.Minor versions
			result = append(result, v)
		} else if v.Segments()[0] != minV.Segments()[0] || v.Segments()[1] != minV.Segments()[1] {
			// Version is newer only within the same Major.Minor version line
			result = append(result, v)
		}
	}
	return result, nil
}

func removeSpecificVersions(versions []*version.Version, specificVersion string) ([]*version.Version, error) {
	specificV, err := version.NewVersion(specificVersion)
	if err != nil {
		return nil, err
	}
	var result []*version.Version
	for _, v := range versions {
		if v.Segments()[0] != specificV.Segments()[0] || v.Segments()[1] != specificV.Segments()[1] || v.Segments()[2] != specificV.Segments()[2] {
			// Version is newer across Major.Minor versions
			result = append(result, v)
		}
	}
	return result, nil
}

func ConvertStringSliceToVersionSlice(strings []string) []*version.Version {

	versions := make([]*version.Version, len(strings))
	for i, raw := range strings {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}

	return versions
}

func ConvertVersionSliceToStringSlice(versions []*version.Version) []string {

	var versionStrings []string
	for _, v := range versions {
		versionStrings = append(versionStrings, v.String())
	}
	return versionStrings
}
