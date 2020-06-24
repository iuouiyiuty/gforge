package dao

import (
	"bytes"
	"io"
	"text/template"
)

const (
	daoCode = `

	func (d *Dao) {{.StructName}}DB(ctx context.Context, id int64) (*{{.StructPkg}}{{.StructName}}, error) {
		db := d.db.Context(ctx)

		var res {{.StructPkg}}{{.StructName}}

		err := db.Table(res.TableName()).Where("id = ?", id).Take(&res).Error
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
	
		return &res, nil
	}

	func (d *Dao) {{.StructName}}GetOne(ctx context.Context, where map[string]interface{}) (*{{.StructPkg}}{{.StructName}}, error) {
		db := d.db.Context(ctx)

		var res {{.StructPkg}}{{.StructName}}

		err := db.Table(res.TableName()).Where(where).Take(&res).Error
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
	
		return &res, nil
	}

	func (d *Dao) {{.StructName}}List(ctx context.Context, where map[string]interface{}) ([]*{{.StructPkg}}{{.StructName}}, error) {
		db := d.db.Context(ctx)

		var res []*{{.StructPkg}}{{.StructName}}

		err := db.Table({{.StructPkg}}{{.StructName}}{}.TableName()).Where(where).Find(&res).Error
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}
	
		return res, nil
	}

	func (d *Dao) {{.StructName}}Insert(ctx context.Context, tx *gorm.DB, row *{{.StructPkg}}{{.StructName}}) (err error) {
		if tx == nil {
			tx = d.db.Context(ctx).Begin()
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
			if err != nil {
				tx.Rollback()
			}
			err = errors.WithStack(tx.Commit().Error)
		}()
	
		if tx.NewRecord(row) {
			err = errors.WithStack(tx.Create(row).Error)
		}
	
		return
	}

	// 更新指定字段，有updated_at的话需要手动指定
	func (d *Dao) {{.StructName}}Update(ctx context.Context, tx *gorm.DB, where map[string]interface{}, update map[string]interface{}) (affected int64, err error) {
		if tx == nil {
			tx = d.db.Context(ctx).Begin()
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
			if err != nil {
				tx.Rollback()
			}
			err = errors.WithStack(tx.Commit().Error)
		}()
	
		db := tx.Table({{.StructPkg}}{{.StructName}}{}.TableName()).Where(where).Updates(update)
		if db.Error != nil {
			err = errors.WithStack(db.Error)
			return
		}
	
		affected = db.RowsAffected
	
		return
	}

	// 慎用，一般都是软删除
	func (d *Dao) {{.StructName}}Delete(ctx context.Context, tx *gorm.DB, where map[string]interface{}) (affected int64, err error) {
		if tx == nil {
			tx = d.db.Context(ctx).Begin()
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
			if err != nil {
				tx.Rollback()
			}
			err = errors.WithStack(tx.Commit().Error)
		}()
	
		db := tx.Where(where).Delete({{.StructPkg}}{{.StructName}}{})
		if db.Error != nil {
			err = errors.WithStack(db.Error)
			return
		}
	
		affected = db.RowsAffected
	
		return
	}
	
	`
)

type fillData struct {
	StructName string
	StructPkg  string
	TableName  string
}

// GenerateDao generates Dao code
func GenerateDao(tableName, structName string, structPkg string) (io.Reader, error) {
	if structPkg != "" {
		// add dot
		structPkg += "."
	}
	var buff bytes.Buffer
	err := template.Must(template.New("dao").Parse(daoCode)).Execute(&buff, fillData{
		StructName: structName,
		StructPkg:  structPkg,
		TableName:  tableName,
	})
	if nil != err {
		return nil, err
	}
	return &buff, nil
}
