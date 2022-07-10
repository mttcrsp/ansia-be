package feeds

type Collections struct {
	Main     []Feed
	Regional []Feed
	Media    []Feed
}

func (c *Collections) All() []Feed {
	var feeds []Feed
	feeds = append(feeds, c.Main...)
	feeds = append(feeds, c.Regional...)
	feeds = append(feeds, c.Media...)
	return feeds
}
