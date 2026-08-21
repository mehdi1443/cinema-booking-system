package dto

type PaginationRequest struct {
	Page  int `form:"page,default=1" binding:"omitempty,gte=1,lte=1000"`
	Limit int `form:"limit,default=10" binding:"omitempty,gte=1,lte=100"`
}

func (p *PaginationRequest) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	} else if p.Page > 1000 {
		p.Page = 1000
	}

	if p.Limit < 1 {
		p.Limit = 10
	} else if p.Limit > 100 {
		p.Limit = 100
	}
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
