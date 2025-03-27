package dao

import (
	"JenArgo/db"
	"JenArgo/model/po"
	"JenArgo/model/vo"
	"errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var Deploy deploy

type deploy struct{}

// List 列表
func (*deploy) List(en, appName, repoName string, page, size int) (*vo.DeploysListResponse, error) {
	startSet := (page - 1) * size

	var (
		deployList = make([]*po.Deploy, 0)
		total      = 0
	)

	tx := db.GORM.Model(&po.Deploy{}).
		Where("app_name like ?", "%"+appName+"%").
		Where("repo_name like ?", "%"+repoName+"%").
		Where("en like ?", "%"+en+"%").
		Count(&total)

	if tx.Error != nil {
		zap.L().Error("获取Deploy列表失败," + tx.Error.Error())
		return nil, errors.New("获取Deploy列表失败," + tx.Error.Error())
	}

	//分页数据
	tx = db.GORM.Model(&po.Deploy{}).
		Where("app_name like ?", "%"+appName+"%").
		Where("repo_name like ?", "%"+repoName+"%").
		Where("en like ?", "%"+en+"%").
		Limit(size).
		Offset(startSet).
		Order("id desc").
		Find(&deployList)

	if tx.Error != nil {
		zap.L().Error("获取Deploy列表失败," + tx.Error.Error())
		return nil, errors.New("获取Deploy列表失败," + tx.Error.Error())
	}

	return &vo.DeploysListResponse{
		Items: deployList,
		Total: total,
	}, nil

}

// Get 查询单个
func (*deploy) Get(deployId int64) (*vo.DeployRequest, bool, error) {
	data := new(vo.DeployRequest)
	tx := db.GORM.Model(&po.Deploy{}).Where("id = ?", deployId).First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if tx.Error != nil {
		zap.L().Error("查询Deploy失败," + tx.Error.Error())
		return nil, false, errors.New("查询Deploy失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Add 新增
func (*deploy) Add(d *po.Deploy) error {
	tx := db.GORM.Create(&d)
	if tx.Error != nil {
		zap.L().Error("新增Deploy失败," + tx.Error.Error())
		return errors.New("新增Deploy失败," + tx.Error.Error())
	}

	return nil
}

// Has 根据应用名查询，用于代码层去重
func (*deploy) Has(en, appName, repoName string) (*po.Deploy, bool, error) {
	data := new(po.Deploy)
	tx := db.GORM.Where("en = ? and app_name = ? and repo_name = ?", en, appName, repoName).Order("created_at desc").First(&data)
	tx.Debug()
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if tx.Error != nil {
		zap.L().Error("根据环境、appName、repoName查询Deploy失败," + tx.Error.Error())
		return nil, false, errors.New("根据环境、appName、repoName查询Deploy失败," + tx.Error.Error())
	}

	return data, true, nil
}

// Update 更新
func (*deploy) Update(d *po.Deploy) error {
	tx := db.GORM.Model(&po.Deploy{}).Where("id = ?", d.ID).Updates(&d)
	//tx := db.GORM.Model(&po.Deploy{}).Where("id = ?", d.ID).Select("tag").Updates(&d)
	if tx.Error != nil {
		zap.L().Error("更新Deploy失败," + tx.Error.Error())
		return errors.New("更新Deploy失败," + tx.Error.Error())
	}
	tx = db.GORM.Model(&po.Deploy{}).Where("id = ?", d.ID).Updates(map[string]interface{}{
		"tag": d.Tag,
	})
	if tx.Error != nil {
		zap.L().Error("更新Deploy失败," + tx.Error.Error())
		return errors.New("更新Deploy失败," + tx.Error.Error())
	}
	if tx.Error != nil {
		zap.L().Error("更新Deploy失败," + tx.Error.Error())
		return errors.New("更新Deploy失败," + tx.Error.Error())
	}

	return nil
}

// Delete 删除
func (*deploy) Delete(deployId int64) error {
	data := new(po.Deploy)
	data.ID = deployId
	tx := db.GORM.Delete(&data)
	if tx.Error != nil {
		zap.L().Error("删除Deploy失败," + tx.Error.Error())
		return errors.New("删除Deploy失败," + tx.Error.Error())
	}

	return nil
}
