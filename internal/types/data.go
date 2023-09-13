package types

import cnst "playground/internal/constants"

// JsonFileData represents the JSON file records structure
type JsonFileData struct {
	CampaignId string  `json:"CampaignId"`
	Country    string  `json:"Country"`
	Ltv1       float64 `json:"Ltv1"`
	Ltv2       float64 `json:"Ltv2"`
	Ltv3       float64 `json:"Ltv3"`
	Ltv4       float64 `json:"Ltv4"`
	Ltv5       float64 `json:"Ltv5"`
	Ltv6       float64 `json:"Ltv6"`
	Ltv7       float64 `json:"Ltv7"`
	Users      int     `json:"Users"`
}

// LtvCollection represents a LTV (Lifetime Value) data
type LtvCollection [cnst.LtvLen]float64

// Record struct represents a common data type retrieved from data sources
// That means that all data sources should provide Record data in system
type Record struct {
	campaignId, country string
	ltv                 LtvCollection
}

// NewRecord initializes and returns a new Record struct
func NewRecord(campaignId, country string, ltv LtvCollection) *Record {
	return &Record{
		campaignId: campaignId,
		country:    country,
		ltv:        ltv,
	}
}

// Record struct getters
func (r *Record) CampaignId() string { return r.campaignId }
func (r *Record) Country() string    { return r.country }
func (r *Record) Ltv() LtvCollection { return r.ltv }

type KeyType string

// AggregatedData struct represents aggregated data, according to key
type AggregatedData struct {
	key KeyType
	ltv LtvCollection
}

// NewAggregatedData initializes and returns a new AggregatedData struct
func NewAggregatedData(key KeyType, ltv LtvCollection) *AggregatedData {
	return &AggregatedData{
		key: key,
		ltv: ltv,
	}
}

// AggregatedData struct getters
func (r *AggregatedData) Key() KeyType       { return r.key }
func (r *AggregatedData) Ltv() LtvCollection { return r.ltv }
