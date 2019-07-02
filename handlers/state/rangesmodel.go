package state

func Single(i *Item) *Ranges {
	return &Ranges{single: true, items: []*Item{i}}
}

func List(is ...*Item) *Ranges {
	return &Ranges{items: is}
}

func Point(key string) *Item {
	return &Item{point: true, begin: key}
}

func Range(begin, until string) *Item {
	return &Item{begin: begin, until: until}
}

type Ranges struct {
	single bool
	items  []*Item
}

func (q *Ranges) IsSingle() bool {
	return q.single
}

func (q *Ranges) Single() *Item {
	return q.items[0]
}

func (q *Ranges) List() []*Item {
	return q.items
}

type Item struct {
	point bool
	begin string
	until string
}

func (i *Item) IsPoint() bool {
	return i.point
}

func (i *Item) Point() string {
	return i.begin
}

func (i *Item) Range() (string, string) {
	return i.begin, i.until
}
