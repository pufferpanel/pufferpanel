package services

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestTemplate_GetImportableTemplates(t1 *testing.T) {
	t1.Run("GetImportableTemplates", func(t1 *testing.T) {
		t := &Template{}

		got, err := t.GetImportableTemplates()
		if err != nil {
			t1.Errorf("GetImportableTemplates() error = %v", err)
			return
		}
		if len(got) == 0 {
			t1.Error("GetImportableTemplates() got no files from repo")
			return
		}
		fmt.Printf("%v", got)
	})
}

func TestTemplate_ImportTemplates(t1 *testing.T) {
	t1.Run("GetImportableTemplates", func(t1 *testing.T) {
		db := prepareDatabase(t1)
		if t1.Failed() {
			return
		}
		t := &Template{DB: db}

		got, err := t.getTemplateFiles()
		if err != nil {
			t1.Errorf("GetImportableTemplates() error = %v", err)
			return
		}
		if len(got) == 0 {
			t1.Error("GetImportableTemplates() got no files from repo")
			return
		}

		for _, template := range got {
			err = t.ImportFromRepo(template.Name)
			if err != nil {
				t1.Errorf("ImportFromRepo() error = %v", err)
				return
			}
		}

		var count int64
		err = db.Find(&models.Template{}).Count(&count).Error
		if err != nil {
			t1.Errorf("countTemplates error = %v", err)
			return
		}
		if len(got) != int(count) {
			t1.Errorf("ImportFromRepo() expected = %d, added %d", len(got), count)
		}

		expected := 0
		for _, te := range got {
			if te.ReadmePath != "" {
				expected++
			}
		}

		if expected == 0 {
			t1.Errorf("countReadmes expected readmes but did not import any")
			return
		}

		err = db.Find(&models.Template{}).Where("readme <> ?", "").Count(&count).Error
		if err != nil {
			t1.Errorf("countReadmes error = %v", err)
			return
		}

		if expected != int(count) {
			t1.Errorf("countReadmes expected readmes = %d, added %d", expected, count)
			return
		}
	})
}

func prepareDatabase(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file:test.db?cache=shared&mode=memory"), &gorm.Config{})
	if err != nil {
		t.Errorf("prepareDatabase() error = %v", err)
		return nil
	}

	err = db.Migrator().CreateTable(&models.Template{})
	if err != nil {
		t.Errorf("prepareDatabase() error = %v", err)
		return nil
	}

	return db
}
