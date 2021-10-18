package items

type Item struct {
	Id          int64       `json:"id"`
	Seller      int64       `json:"seller"`
	Title       string      `json:"title"`
	Description Description `json:"description"`
	Pictures    []Picture   `json:"pictures"`
	Price       float32     `json:"price"`
}

type Description struct {
	PlainText string `json:"plain_text"`
	Html      string `json:"html"`
}

type Picture struct {
	Id  int64  `json:"id"`
	Url string `json:"url"`
}

func (i *Item) GetPictures(id int64) [][]interface{} {
	var FormattedPics [][]interface{}
	for _, p := range i.Pictures {
		FormattedPics = append(FormattedPics, []interface{}{id, p.Url})
	}
	return FormattedPics
}
