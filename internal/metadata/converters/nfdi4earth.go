package converters

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/alexander-lindner/go-cff"
	toolspec "github.com/hydrocode-de/tool-spec-go"
)

type NFDI4Earth struct {
	Context             string             `json:"@context"`
	ID                  string             `json:"@id,omitempty"`
	Type                string             `json:"@type"`
	N4EID               string             `json:"n4e:id"`
	Name                string             `json:"name"`
	Description         string             `json:"description,omitempty"`
	Version             string             `json:"version,omitempty"`
	Author              []NFDI4EarthPerson `json:"author,omitempty"`
	Contributor         []NFDI4EarthPerson `json:"contributor,omitempty"`
	License             string             `json:"license,omitempty"`
	CodeRepository      string             `json:"codeRepository,omitempty"`
	Keywords            []string           `json:"keywords,omitempty"`
	ProgrammingLanguage []string           `json:"programmingLanguage,omitempty"`
	DateModified        string             `json:"dateModified,omitempty"`
	DatePublished       string             `json:"datePublished,omitempty"`
	Publisher           []NFDI4EarthOrg    `json:"publisher,omitempty"`
	SameAs              []string           `json:"sameAs,omitempty"`
	Identifier          []string           `json:"identifier,omitempty"`
	CopyrightNotice     string             `json:"copyrightNotice,omitempty"`
}

type NFDI4EarthPerson struct {
	Type        string         `json:"@type"`
	Name        string         `json:"name"`
	Email       string         `json:"email,omitempty"`
	OrcidID     string         `json:"orcidId,omitempty"`
	Affiliation *NFDI4EarthOrg `json:"affiliation,omitempty"`
	URL         string         `json:"url,omitempty"`
}

type NFDI4EarthOrg struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type NFDI4EarthConverter struct {
	NFDI4Earth
	errs []error
}

func (n *NFDI4EarthConverter) Ingest(spec toolspec.ToolSpec) {
	n.NFDI4Earth = NFDI4Earth{
		Context: "https://nfdi4earth.pages.rwth-aachen.de/knowledgehub/nfdi4earth-kh-schema/SoftwareSourceCode/context.jsonld",
		Type:    "n4e:SoftwareSourceCode",
		Name:    spec.Title,
	}

	if spec.Description != "" {
		n.NFDI4Earth.Description = spec.Description
	}

	if spec.Citation.Title != "" {
		if spec.Citation.Version != "" {
			n.NFDI4Earth.Version = spec.Citation.Version
		}

		authors := []string{}
		for _, author := range spec.Citation.Authors {
			if author.IsPerson {
				fullName := strings.TrimSpace(author.Person.GivenNames + " " + author.Person.Family)
				authors = append(authors, fullName)
				p := NFDI4EarthPerson{
					Type: "n4e:Person",
					Name: fullName,
				}
				if author.Person.Email != "" {
					p.Email = author.Person.Email
				}
				if author.Person.Orcid != "" {
					p.OrcidID = string(author.Person.Orcid)
					if strings.Contains(p.OrcidID, "orcid.org/") {
						u, err := url.Parse(p.OrcidID)
						if err == nil {
							p.OrcidID = strings.TrimPrefix(u.Path, "/")
						}
					}
				}
				if author.Person.Affiliation != "" {
					p.Affiliation = &NFDI4EarthOrg{
						Type: "n4e:Organization",
						Name: author.Person.Affiliation,
					}
				}
				if author.Person.Website.URL != nil {
					p.URL = author.Person.Website.URL.String()
				}
				n.NFDI4Earth.Author = append(n.NFDI4Earth.Author, p)
			} else {
				authors = append(authors, author.Entity.Name)
				n.NFDI4Earth.Publisher = append(n.NFDI4Earth.Publisher, NFDI4EarthOrg{
					Type: "n4e:Organization",
					Name: author.Entity.Name,
					URL:  n.getUrlString(author.Entity.Website),
				})
			}
		}

		if repo := n.getUrlString(spec.Citation.RepositoryCode); repo != "" {
			n.NFDI4Earth.CodeRepository = repo
			n.NFDI4Earth.N4EID = repo
		} else if repo := n.getUrlString(spec.Citation.Url); repo != "" {
			n.NFDI4Earth.CodeRepository = repo
			n.NFDI4Earth.N4EID = repo
		}

		t := time.Time(spec.Citation.DateReleased)
		if !t.IsZero() {
			n.NFDI4Earth.DatePublished = t.Format("2006-01-02")
			n.NFDI4Earth.CopyrightNotice = fmt.Sprintf("Copyright %d, %s", t.Year(), strings.Join(authors, ", "))
		} else {
			n.NFDI4Earth.CopyrightNotice = fmt.Sprintf("Copyright %d, %s", time.Now().Year(), strings.Join(authors, ", "))
		}

		if spec.Citation.Doi.General != "" {
			doi := fmt.Sprintf("%s.%s/%s", spec.Citation.Doi.General, spec.Citation.Doi.DirectoryIndicator, spec.Citation.Doi.RegistrantCode)
			doiUrl := doi
			if !strings.HasPrefix(doiUrl, "http") {
				doiUrl = "https://doi.org/" + doi
			}
			n.NFDI4Earth.Identifier = append(n.NFDI4Earth.Identifier, doi)
			n.NFDI4Earth.SameAs = append(n.NFDI4Earth.SameAs, doiUrl)
			n.NFDI4Earth.ID = doiUrl
			n.NFDI4Earth.N4EID = doi
		}

		if len(spec.Citation.License.Data) > 0 {
			lic := string(spec.Citation.License.Data[0])
			if !strings.HasPrefix(lic, "http") {
				lic = "http://spdx.org/licenses/" + lic
			}
			n.NFDI4Earth.License = lic
		}

		n.NFDI4Earth.Keywords = spec.Citation.Keywords
	}

	if n.N4EID == "" {
		n.N4EID = strings.ToLower(strings.ReplaceAll(spec.Title, " ", "-"))
	}
}

func (n *NFDI4EarthConverter) getUrlString(u cff.URL) string {
	if u.URL != nil {
		return u.URL.String()
	}
	return ""
}

func (n *NFDI4EarthConverter) Validate() bool {
	if n.Name == "" {
		n.errs = append(n.errs, fmt.Errorf("name is required"))
	}
	if n.N4EID == "" {
		n.errs = append(n.errs, fmt.Errorf("n4e:id is required"))
	}
	return len(n.errs) == 0
}

func (n *NFDI4EarthConverter) Serialize(format string) ([]byte, error) {
	return json.MarshalIndent(n.NFDI4Earth, "", "  ")
}
