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
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BeerBox struct {
	ID    int64 `json:"id"`
	Price int64 `json:"price"`
}
