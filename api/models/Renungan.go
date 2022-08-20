package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Renungan struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Renungan) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Renungan) Validate() error {

	if p.Title == "" {
		return errors.New("required Title")
	}
	if p.Content == "" {
		return errors.New("required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("required Author")
	}
	return nil
}

func (p *Renungan) SaveRenungan(db *gorm.DB) (*Renungan, error) {
	var err error
	err = db.Debug().Model(&Renungan{}).Create(&p).Error
	if err != nil {
		return &Renungan{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Renungan{}, err
		}
	}
	return p, nil
}

func (p *Renungan) FindAllRenungan(db *gorm.DB) (*[]Renungan, error) {
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

func (p *Renungan) FindRenunganByID(db *gorm.DB, pid uint64) (*Renungan, error) {
	var err error
	err = db.Debug().Model(&Renungan{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Renungan{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Renungan{}, err
		}
	}
	return p, nil
}

func (p *Renungan) UpdateARenungan(db *gorm.DB) (*Renungan, error) {

	var err error

	err = db.Debug().Model(&Renungan{}).Where("id = ?", p.ID).Updates(Renungan{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Renungan{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Renungan{}, err
		}
	}
	return p, nil
}

func (p *Renungan) DeleteARenungan(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Renungan{}).Where("id = ? and author_id = ?", pid, uid).Take(&Renungan{}).Delete(&Renungan{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}