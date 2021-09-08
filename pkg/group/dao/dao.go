package dao

import (
	"context"

	"g.hz.netease.com/horizon/lib/orm"
	"g.hz.netease.com/horizon/pkg/group/models"
)

type DAO interface {
	Create(ctx context.Context, group *models.Group) (int64, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*models.Group, error)
	GetByPath(ctx context.Context, path string) (*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	GetChildren(ctx context.Context, id int64) ([]*models.Group, error)
}

// New returns an instance of the default DAO
func New() DAO {
	return &dao{}
}

type dao struct{}

func (d *dao) Create(ctx context.Context, group *models.Group) (int64, error) {
	// todo 判断同名
	db, err := orm.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	result := db.Create(group)

	return result.RowsAffected, result.Error
}

func (d *dao) Delete(ctx context.Context, id int64) error {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}

	result := db.Where("id = ?", id).Delete(&models.Group{})

	return result.Error
}

func (d *dao) Get(ctx context.Context, id int64) (*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var group *models.Group
	result := db.First(&group, id)

	return group, result.Error
}

func (d *dao) GetByPath(ctx context.Context, path string) (*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var group *models.Group
	result := db.Find(&group, "path = ?", path)

	return group, result.Error
}

func (d *dao) GetChildren(ctx context.Context, id int64) ([]*models.Group, error) {
	db, err := orm.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group
	result := db.Where("parentId = ?", id).Find(&groups)

	return groups, result.Error
}

func (d *dao) Update(ctx context.Context, group *models.Group) error {
	// todo 判断同名
	db, err := orm.FromContext(ctx)
	if err != nil {
		return err
	}

	result := db.Model(group).Select("Name", "Description", "VisibilityLevel").Updates(group)

	return result.Error
}
