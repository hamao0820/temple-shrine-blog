package util

type Paginate struct {
	BeforePage int
	NextPage   int
}

func CreatePaginate(page int, limit int, count int) Paginate {
	beforePage := page - 1
	if beforePage <= 0 {
		beforePage = 0
	}

	nextPage := page + 1
	if count != limit {
		nextPage = 0
	}

	return Paginate{
		BeforePage: beforePage,
		NextPage:   nextPage,
	}
}
