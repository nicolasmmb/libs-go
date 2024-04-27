package paginator

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	MAX_ITEMS_PER_PAGE     = ""
	MIN_ITEMS_PER_PAGE     = ""
	MAX_ITEMS_PER_PAGE_INT = 0
	MIN_ITEMS_PER_PAGE_INT = 0

	ErrPageMustBeANumber  = errors.New("page must be a number")
	ErrLimitMustBeANumber = errors.New("limit must be a number")
	ErrLimitMin           = errors.New("limit min is: ")
	ErrLimitMax           = errors.New("limit max is: ")
)

func init() {
	log.Println("Inicializando: helper/pagination")
	MAX_ITEMS_PER_PAGE = os.Getenv("MAX_ITEMS_PER_PAGE")
	MIN_ITEMS_PER_PAGE = os.Getenv("MIN_ITEMS_PER_PAGE")

	_MAX_ITEMS_PER_PAGE_INT, err := strconv.Atoi(MAX_ITEMS_PER_PAGE)
	if err != nil {
		log.Println("Erro ao converter MAX_ITEMS_PER_PAGE para inteiro, utilizando valor padrão '100'")
		_MAX_ITEMS_PER_PAGE_INT = 100
	}
	_MIN_ITEMS_PER_PAGE_INT, err := strconv.Atoi(MIN_ITEMS_PER_PAGE)
	if err != nil {
		log.Println("Erro ao converter MIN_ITEMS_PER_PAGE para inteiro, utilizando valor padrão '1'")
		_MIN_ITEMS_PER_PAGE_INT = 1
	}

	MAX_ITEMS_PER_PAGE_INT = _MAX_ITEMS_PER_PAGE_INT
	MIN_ITEMS_PER_PAGE_INT = _MIN_ITEMS_PER_PAGE_INT

	if MAX_ITEMS_PER_PAGE_INT < MIN_ITEMS_PER_PAGE_INT {
		log.Println("MAX_ITEMS_PER_PAGE não pode ser menor que MIN_ITEMS_PER_PAGE, utilizando valor padrão '100' para MAX_ITEMS_PER_PAGE e '1' para MIN_ITEMS_PER_PAGE")
		MAX_ITEMS_PER_PAGE_INT = 100
		MIN_ITEMS_PER_PAGE_INT = 1
	}

	ErrLimitMin = errors.New("limit min is: " + MIN_ITEMS_PER_PAGE)
	ErrLimitMax = errors.New("limit max is: " + MAX_ITEMS_PER_PAGE)
	log.Println("MAX_ITEMS_PER_PAGE: ", MAX_ITEMS_PER_PAGE_INT)
	log.Println("MIN_ITEMS_PER_PAGE: ", MIN_ITEMS_PER_PAGE_INT)
}

type Pagination struct {
	Page  int
	Limit int
}

func Create(c *gin.Context) (pagination *Pagination, err error) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, ErrPageMustBeANumber
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, ErrLimitMustBeANumber
	}

	if limit < MIN_ITEMS_PER_PAGE_INT {
		return nil, ErrLimitMin
	}
	if limit > MAX_ITEMS_PER_PAGE_INT {
		return nil, ErrLimitMax
	}

	return &Pagination{Page: page, Limit: limit}, nil
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
