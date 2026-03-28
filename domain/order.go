package domain

type Tag string

const (
	Tag315Pro  Tag = "315pro"
	TagGakumas Tag = "学マス"
	Tag283Pro  Tag = "283pro"
	Tag346Pro  Tag = "346pro"
	Tag765Pro  Tag = "765pro"
)

// ValidTags is the authoritative list of allowed tag values.
var ValidTags = []Tag{Tag315Pro, TagGakumas, Tag283Pro, Tag346Pro, Tag765Pro}

type Order struct {
	ThreadName string
	Deadline   string // ISO-8601 date, may be empty
	Tag        Tag    // single select, may be empty
	ShopURL    string // not persisted to Notion
}
