package eagle

type Pack struct {
	Images []Image `json:"images"`
}

type Image struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Size             int64    `json:"size"`
	BTime            int64    `json:"btime"`
	MTime            int64    `json:"mtime"`
	Ext              string   `json:"ext"`
	Tags             []string `json:"tags"`
	Folders          []string `json:"folders"`
	IsDeleted        bool     `json:"isDeleted"`
	URL              string   `json:"url"`
	Annotation       string   `json:"annotation"`
	ModificationTime int64    `json:"modificationTime"`
	Star             int      `json:"star,omitempty"`
	NoThumbnail      bool     `json:"noThumbnail"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
	Palettes         []string `json:"palettes"`
	LastModified     int64    `json:"lastModified"`
}
