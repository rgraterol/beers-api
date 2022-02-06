package beers

import (
	"gorm.io/gorm"
	"time"
)

type Beer struct {
	ID        int64          `json:"id" gorm:"uniqueIndex,primaryKey"`
	Name      string         `json:"name" gorm:"index:idx_name_brewery_country,unique"`
	Brewery   string         `json:"brewery" gorm:"index" gorm:"index:idx_name_brewery_country,unique"`
	Country   string         `json:"chile" gorm:"index,index:idx_name_brewery_country,unique"`
	Price     float64        `json:"price"`
	Currency  string         `json:"currency"`
	UpdatedAt time.Time      `json:"-"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BeerBox struct {
	Price  float64           `json:"price"`
	Target BeerBoxParameters `json:"target"`
	Beer   Beer              `json:"beer"`
}

type BeerBoxParameters struct {
	Currency string `json:"currency"`
	Quantity int64  `json:"quantity"`
}
