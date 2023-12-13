package utils

import (
	"github.com/elastic/go-elasticsearch/v8"
	"tf_ocg/proto/models"
)

var EsClient *elasticsearch.Client

type ProductSearchResult struct {
	Products   []*models.Product
	TotalItems int64
}
