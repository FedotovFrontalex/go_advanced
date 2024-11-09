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
