package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Renungan struct {
	gorm.Model
	Title    string `gorm:"size:255;not null;unique" json:"title"`
	Content  string `gorm:"size:255;not null;" json:"content"`
	Author   User   `json:"author"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
}

func (r *Renungan) Prepare() {
	r.ID = 0
	r.Title = html.EscapeString(strings.TrimSpace(r.Title))
	r.Content = html.EscapeString(strings.TrimSpace(r.Content))
	r.Author = User{}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}

func (r *Renungan) Validate() error {

	if r.Title == "" {
		return errors.New("required Title")
	}
	if r.Content == "" {
		return errors.New("required Content")
	}
	if r.AuthorID < 1 {
		return errors.New("required Author")
	}
	return nil
}

func (r *Renungan) SaveRenungan(db *gorm.DB) (*Renungan, error) {
	var err error
	err = db.Debug().Model(&Renungan{}).Create(&r).Error
	if err != nil {
		return &Renungan{}, err
	}
	if r.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", r.AuthorID).Take(&r.Author).Error
		if err != nil {
			return &Renungan{}, err
		}
	}
	return r, nil
}

func (r *Renungan) FindAllRenungan(db *gorm.DB) (*[]Renungan, error) {
	var err error
	posts := []Renungan{}
	err = db.Debug().Model(&Renungan{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Renungan{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Renungan{}, err
			}
		}
	}
	return &posts, nil
}

func (r *Renungan) FindRenunganByID(db *gorm.DB, pid uint) (*Renungan, error) {
	var err error
	err = db.Debug().Model(&Renungan{}).Where("id = ?", pid).Take(&r).Error
	if err != nil {
		return &Renungan{}, err
	}
	if r.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", r.AuthorID).Take(&r.Author).Error
		if err != nil {
			return &Renungan{}, err
		}
	}
	return r, nil
}

func (r *Renungan) UpdateARenungan(db *gorm.DB) (*Renungan, error) {

	var err error

	err = db.Debug().Model(&Renungan{}).Where("id = ?", r.ID).Updates(Renungan{Title: r.Title, Content: r.Content}).Error
	if err != nil {
		return &Renungan{}, err
	}
	if r.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", r.AuthorID).Take(&r.Author).Error
		if err != nil {
			return &Renungan{}, err
		}
	}
	return r, nil
}

func (r *Renungan) DeleteARenungan(db *gorm.DB, pid uint, uid uint) (int64, error) {

	db = db.Debug().Model(&Renungan{}).Where("id = ? and author_id = ?", pid, uid).Take(&Renungan{}).Delete(&Renungan{})

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
