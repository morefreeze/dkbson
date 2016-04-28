package duokan

// Bound each item has a bound.
type Bound struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

// Item each item is a word or a picture in a page.
// In fact, we don't need X,Y instead Bound.X/Y.
type Item struct {
	Bound Bound  `json:"bound"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Pos   []int  `json:"pos"`
	Type  string `json:"type"`
	Char  string `json:"char"`
}

var paragraphStart = "    "

// PageContent contains many items which represent element in page.
type PageContent struct {
	Items []*Item
}

// GenerateContent generates content.
func (p *PageContent) GenerateContent() (string, error) {
	if len(p.Items) <= 0 {
		return "", nil
	}
	// paragraph left, it is the most left which is almost X > left.
	left, up := 10000, 10000
	for _, item := range p.Items {
		if item.Type != "word" {
			continue
		}
		if item.Bound.X < left {
			left = item.Bound.X
		}
		if item.Bound.Y < up {
			up = item.Bound.Y
		}
	}
	content := ""
	lastY := up
	for _, item := range p.Items {
		switch item.Type {
		case "word":
			if item.Bound.X > left && lastY != item.Bound.Y {
				content += "\n    "
			}
			lastY = item.Bound.Y
			content += item.Char
		case "image":
			// TODO: handle picture?
		}
	}
	return content, nil
}
