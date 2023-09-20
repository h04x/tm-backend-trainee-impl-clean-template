// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "encoding/json"

// Metrics -.
type Metrics struct {
	Date   string `json:"date"         binding:"required,datetime=2006-01-02" example:"2006-01-02"`
	Views  uint   `json:"views"                                               example:"5"`
	Clicks uint   `json:"clicks"                                              example:"5"`
	Cost   string `json:"cost"         binding:"ValidateCost"                 example:"1.25"`
}

// Statistics -.
type Statistics struct {
	Metrics Metrics
	Cpc     float32 `json:"cpc"        example:"0.21"`
	Cpm     float32 `json:"cpm"        example:"1.23"`
}

// DoGetRequest -.
type DoGetRequest struct {
	From  string `json:"from"          binding:"required,datetime=2006-01-02"             example:"2006-01-02" :"from"`
	To    string `json:"to"            binding:"required,datetime=2006-01-02"             example:"2006-01-02" :"to"`
	Order string `json:"order"         binding:"oneof=Date Views Clicks Cost Cpc Cpm"     example:"cpc" :"order"`
}

func (s *Statistics) MarshalJSON() ([]byte, error) {
	intermediate := map[string]interface{}{
		"date":   s.Metrics.Date,
		"views":  s.Metrics.Views,
		"clicks": s.Metrics.Clicks,
		"cost":   s.Metrics.Cost,
		"cpc":    s.Cpc,
		"cpm":    s.Cpm,
	}
	return json.Marshal(intermediate)
}
