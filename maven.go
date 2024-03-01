package installer

import (
	"encoding/xml"
	"net/url"
)

type (
	MavenMetadataVersioning struct {
		Release     string   `xml:"release"`
		Latest      string   `xml:"latest"`
		LastUpdated string   `xml:"lastUpdated"`
		Versions    []string `xml:"versions>version"`
	}

	MavenMetadata struct {
		GroupId    string                  `xml:"groupId"`
		ArtifactId string                  `xml:"artifactId"`
		Versioning MavenMetadataVersioning `xml:"versioning"`
	}
)

func DecodeMavenMetadata(body []byte) (data MavenMetadata, err error) {
	if err = xml.Unmarshal(body, &data); err != nil {
		return
	}
	return
}

func GetMavenMetadata(link string) (data MavenMetadata, err error) {
	if link, err = url.JoinPath(link, "maven-metadata.xml"); err != nil {
		return
	}
	if err = DefaultHTTPClient.GetXml(link, &data); err != nil {
		return
	}
	return
}
