package types

import cnst "playground/internal/constants"

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
