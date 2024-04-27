package paginator

type PaginationResponse[R any] struct {
	Msg        string `json:"msg"`
	Total      *int   `json:"total"`
	Limit      *int   `json:"limit"`
	PageLast   *int   `json:"page_last"`
	PageActual *int   `json:"page_actual"`
	ItensTotal int    `json:"itens_total"`
	Itens      []*R   `json:"itens"`
}

func CreatePaginationResponse[R any](msg string, page *Pagination, total *int, itens []*R) *PaginationResponse[R] {
	p := &PaginationResponse[R]{
		Msg:        msg,
		Total:      total,
		Limit:      &page.Limit,
		PageActual: &page.Page,
		Itens:      itens,
		ItensTotal: len(itens),
	}
	if p.Total != nil {
		p.UpdateLastPage()
	}

	return p
}

func (p *PaginationResponse[R]) UpdateLastPage() {
	total := *p.Total / *p.Limit
	if *p.Total%*p.Limit != 0 {
		total++
	}
	p.PageLast = &total
}
