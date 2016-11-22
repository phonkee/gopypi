package core

import "time"

type StatsAggregation int

const (
	STATS_DOWNLOAD_WEEKLY StatsAggregation = iota + 1
	STATS_DOWNLOAD_MONTHLY
	STATS_DOWNLOAD_ALL
)

type StatsDownloadItem struct {
	PackageVersion   *PackageVersion `gorm:"ForeignKey:PackageVersionID" json:"version,omitempty"`
	PackageVersionID uint            `json:"-"`
	Downloads        int             `json:"downloads"`
	CreatedAt        time.Time       `json:"created_at"`
}
