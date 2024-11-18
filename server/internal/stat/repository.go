package stat

import (
	"server/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())

	repo.DB.Find(&stat, "link_id=? and date=?", linkId, currentDate)

	if stat.ID == 0 {
		repo.DB.Create(&Stat{
			LinkId: linkId,
			Date:   currentDate,
			Clicks: 1,
		})
	} else {
		stat.Clicks += 1
		repo.DB.Save(stat)
	}
}

func (repo *StatRepository) GetStats(by string, from time.Time, to time.Time) []StatGetResponse {
	var stats []StatGetResponse
	var selectQuery string

	switch by {
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	}

	repo.DB.Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats
}
