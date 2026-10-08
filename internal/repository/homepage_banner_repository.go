package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type HomepageBannerRepository struct {
	db *gorm.DB
}

func NewHomepageBannerRepository(db *gorm.DB) *HomepageBannerRepository {
	return &HomepageBannerRepository{db: db}
}

// FindAll mengambil seluruh slide HomepageBanner (aktif maupun non-aktif) beserta buttons
func (r *HomepageBannerRepository) FindAll() ([]model.HomepageBanner, error) {
	var banners []model.HomepageBanner
	err := r.db.
		Preload("Buttons", func(db *gorm.DB) *gorm.DB {
			return db.Order("homepage_banner_buttons.order ASC, homepage_banner_buttons.id ASC")
		}).
		Order("homepage_banners.order ASC, homepage_banners.id ASC").
		Find(&banners).Error
	if err != nil {
		return nil, err
	}
	for i := range banners {
		if banners[i].Buttons == nil {
			banners[i].Buttons = []model.HomepageBannerButton{}
		}
	}
	return banners, nil
}

// FindActive mengambil seluruh slide HomepageBanner yang aktif (is_active = true) beserta buttons
func (r *HomepageBannerRepository) FindActive() ([]model.HomepageBanner, error) {
	var banners []model.HomepageBanner
	err := r.db.Where("is_active = ?", true).
		Preload("Buttons", func(db *gorm.DB) *gorm.DB {
			return db.Order("homepage_banner_buttons.order ASC, homepage_banner_buttons.id ASC")
		}).
		Order("homepage_banners.order ASC, homepage_banners.id ASC").
		Find(&banners).Error
	if err != nil {
		return nil, err
	}
	for i := range banners {
		if banners[i].Buttons == nil {
			banners[i].Buttons = []model.HomepageBannerButton{}
		}
	}
	return banners, nil
}

// FindByID mencari HomepageBanner berdasarkan ID beserta buttons
// OrderExists mengecek apakah nomor urut sudah dipakai banner lain.
// excludeID dipakai saat edit agar banner yang sedang diedit tidak dianggap duplikat.
func (r *HomepageBannerRepository) OrderExists(order int, excludeID uint) (bool, error) {
	query := r.db.Model(&model.HomepageBanner{}).Where("\"order\" = ?", order)

	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *HomepageBannerRepository) FindByID(id uint) (*model.HomepageBanner, error) {
	var banner model.HomepageBanner
	err := r.db.
		Preload("Buttons", func(db *gorm.DB) *gorm.DB {
			return db.Order("homepage_banner_buttons.order ASC, homepage_banner_buttons.id ASC")
		}).
		First(&banner, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if banner.Buttons == nil {
		banner.Buttons = []model.HomepageBannerButton{}
	}
	return &banner, nil
}

// Create membuat HomepageBanner baru beserta relasi HomepageBannerButton dalam satu transaksi DB
func (r *HomepageBannerRepository) Create(banner *model.HomepageBanner, buttons []model.HomepageBannerButton) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		banner.Buttons = nil // Biarkan penanganan child buttons dilakukan secara transaksional
		if err := tx.Create(banner).Error; err != nil {
			return err
		}

		if len(buttons) > 0 {
			for i := range buttons {
				buttons[i].BannerID = banner.ID
			}
			if err := tx.Create(&buttons).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Update memperbarui data HomepageBanner dan mereplace seluruh button miliknya dalam satu transaksi DB
func (r *HomepageBannerRepository) Update(banner *model.HomepageBanner, buttons []model.HomepageBannerButton) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(banner).Error; err != nil {
			return err
		}

		// Hapus button lama milik Banner ini
		if err := tx.Where("banner_id = ?", banner.ID).Delete(&model.HomepageBannerButton{}).Error; err != nil {
			return err
		}

		// Buat button baru jika ada
		if len(buttons) > 0 {
			for i := range buttons {
				buttons[i].ID = 0 // Pastikan ID baru agar dibuat ulang secara bersih
				buttons[i].BannerID = banner.ID
			}
			if err := tx.Create(&buttons).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete menghapus HomepageBanner berdasarkan ID (relasi button akan terhapus via CASCADE)
func (r *HomepageBannerRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("banner_id = ?", id).Delete(&model.HomepageBannerButton{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.HomepageBanner{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
