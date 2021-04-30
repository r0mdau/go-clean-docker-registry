package filter

import (
	"fmt"
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
			imageTagsToDelete = append(imageTagsToDelete, imageTag)
			return imageTagsToDelete, nil
		}

		if ok, _ := regexp.MatchString("^"+imageTag, tag); ok {
			identifier := tag[len(imageBaseName):]

			v, err := version.NewVersion(identifier)
			if err != nil {
				fmt.Println(imageTag + " " + tag + " : not a good version number from * wildcard")
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
