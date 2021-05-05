package filter

import (
	"github.com/hashicorp/go-version"
	"regexp"
	"sort"
)

func MatchAndSortImageTags(tags []string, imageTag string) ([]string, error) {
	var imageTagsToDelete []string
	// todo work with this kind of wildcard : *-master
	tagWildcard := imageTag[len(imageTag)-1:]

	mapping := make(map[string]string)

	imageBaseName := imageTag[:len(imageTag)-1]
	var versions []*version.Version

	for _, tag := range tags {
		if tag == imageTag && tagWildcard != "*" {
			imageTagsToDelete = append(imageTagsToDelete, tag)
			return imageTagsToDelete, nil
		}

		if ok, _ := regexp.MatchString("^"+imageTag, tag); ok {
			identifier := tag[len(imageBaseName):]

			v, err := version.NewVersion(identifier)
			if err != nil {
				// appending matching tags not semver(sioned)
				imageTagsToDelete = append(imageTagsToDelete, tag)
				continue
			}
			versions = append(versions, v)
			mapping[v.String()] = tag
		}
	}

	// sort versions from oldest to newest
	sort.Sort(version.Collection(versions))
	for _, value := range versions {
		imageTagsToDelete = append(imageTagsToDelete, mapping[value.String()])
	}
	return imageTagsToDelete, nil
}
