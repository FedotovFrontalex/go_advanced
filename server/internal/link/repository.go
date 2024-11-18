package link

import (
	"server/pkg/db"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link)

	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}

func (repo *LinkRepository) Get(hash string) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, "hash = ?", hash)

	if result.Error != nil {
		return nil, result.Error
	}

	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Database.Updates(link)

	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}

func (repo *LinkRepository) Delete(id uint64) error {
	result := repo.Database.Delete(&Link{}, "id=?", id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *LinkRepository) checkIsExist(id uint64) (*Link, error) {
	var link Link
	result := repo.Database.First(&link, "id=?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &link, nil
}

func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Count(&count)

	return count
}

func (repo *LinkRepository) GetAll(limit int, offset int) []Link {
	var links []Link

	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Scan(&links)

	return links
}
