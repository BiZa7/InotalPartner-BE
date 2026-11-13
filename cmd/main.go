package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"inotal-be/config"
	"inotal-be/internal/database"
	"inotal-be/internal/repository"
	"inotal-be/internal/service"
	"inotal-be/routes"
)

func main() {
	// 1. Load config dari .env
	config.Load()

	// 2. Koneksi ke PostgreSQL + AutoMigrate
	database.Connect()

	// 3. Seed default roles & dependency injection
	roleRepo := repository.NewRoleRepository(database.DB)
	if err := roleRepo.SeedDefaultRoles(); err != nil {
		log.Fatalf("[seed] Failed to seed default roles: %v", err)
	}

	userRepo := repository.NewUserRepository(database.DB)
	authSvc := service.NewAuthService(userRepo, roleRepo)
	userSvc := service.NewUserService(userRepo, roleRepo)
	roleSvc := service.NewRoleService(roleRepo)

	categoryRepo := repository.NewCategoryRepository(database.DB)
	if err := categoryRepo.SyncArticleCounts(); err != nil {
		log.Printf("[sync] Warning: Gagal sinkronisasi article_count kategori: %v", err)
	}
	categorySvc := service.NewCategoryService(categoryRepo)

	tagRepo := repository.NewTagRepository(database.DB)
	tagSvc := service.NewTagService(tagRepo)

	articleRepo := repository.NewArticleRepository(database.DB)
	articleSvc := service.NewArticleService(articleRepo, categoryRepo, tagRepo)

	homepageBannerRepo := repository.NewHomepageBannerRepository(database.DB)
	homepageBannerSvc := service.NewHomepageBannerService(homepageBannerRepo)

	testimonialRepo := repository.NewTestimonialRepository(database.DB)
	testimonialSvc := service.NewTestimonialService(testimonialRepo)

	aboutSettingRepo := repository.NewAboutPageSettingRepository(database.DB)
	aboutSettingSvc := service.NewAboutPageSettingService(aboutSettingRepo)

	teamMemberRepo := repository.NewTeamMemberRepository(database.DB)
	teamMemberSvc := service.NewTeamMemberService(teamMemberRepo)

	aboutVideoRepo := repository.NewAboutVideoRepository(database.DB)
	aboutVideoSvc := service.NewAboutVideoService(aboutVideoRepo)

	aboutFAQRepo := repository.NewAboutFAQRepository(database.DB)
	aboutFAQSvc := service.NewAboutFAQService(aboutFAQRepo)

	aboutVideoSettingRepo := repository.NewAboutVideoSettingRepository(database.DB)
	aboutVideoSettingSvc := service.NewAboutVideoSettingService(aboutVideoSettingRepo)

	aboutVideoHighlightRepo := repository.NewAboutVideoHighlightRepository(database.DB)
	aboutVideoHighlightSvc := service.NewAboutVideoHighlightService(aboutVideoHighlightRepo)

	aboutStatisticRepo := repository.NewAboutStatisticRepository(database.DB)
	aboutStatisticSvc := service.NewAboutStatisticService(aboutStatisticRepo)

	pageBannerRepo := repository.NewPageBannerRepository(database.DB)
	pageBannerSvc := service.NewPageBannerService(pageBannerRepo)

	homepageReasonRepo := repository.NewHomepageReasonRepository(database.DB)
	homepageReasonSvc := service.NewHomepageReasonService(homepageReasonRepo)

	homeCTARepo := repository.NewHomeCTARepository(database.DB)
	homeCTASvc := service.NewHomeCTAService(homeCTARepo)

	partnershipSettingRepo := repository.NewPartnershipSettingRepository(database.DB)
	partnershipSettingSvc := service.NewPartnershipSettingService(partnershipSettingRepo)

	partnershipImageRepo := repository.NewPartnershipImageRepository(database.DB)
	partnershipImageSvc := service.NewPartnershipImageService(partnershipImageRepo)

	partnershipServiceRepo := repository.NewPartnershipServiceRepository(database.DB)
	partnershipServiceSvc := service.NewPartnershipServiceService(partnershipServiceRepo)

	aboutIntroSettingRepo := repository.NewAboutIntroSettingRepository(database.DB)
	aboutIntroSettingSvc := service.NewAboutIntroSettingService(aboutIntroSettingRepo)

	aboutIntroMetricRepo := repository.NewAboutIntroMetricRepository(database.DB)
	aboutIntroMetricSvc := service.NewAboutIntroMetricService(aboutIntroMetricRepo, aboutIntroSettingRepo)

	aboutIntroPointRepo := repository.NewAboutIntroPointRepository(database.DB)
	aboutIntroPointSvc := service.NewAboutIntroPointService(aboutIntroPointRepo, aboutIntroSettingRepo)

	homepageServicesDocSvc := service.NewHomepageServicesUpsertService(
		database.DB,
		partnershipSettingRepo,
		partnershipImageRepo,
		partnershipServiceRepo,
	)

	aboutIntroDocSvc := service.NewAboutIntroUpsertService(
		database.DB,
		aboutIntroSettingRepo,
		aboutIntroMetricRepo,
		aboutIntroPointRepo,
	)

	// 4. Setup Gin
	if config.App.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 5. Daftarkan semua route
	routes.Setup(
		r,
		authSvc,
		userSvc,
		roleSvc,
		categorySvc,
		tagSvc,
		articleSvc,
		homepageBannerSvc,
		testimonialSvc,
		aboutSettingSvc,
		teamMemberSvc,
		aboutVideoSvc,
		aboutFAQSvc,
		aboutVideoSettingSvc,
		aboutVideoHighlightSvc,
		aboutStatisticSvc,
		pageBannerSvc,
		homepageReasonSvc,
		homeCTASvc,
		partnershipSettingSvc,
		partnershipImageSvc,
		partnershipServiceSvc,
		aboutIntroSettingSvc,
		aboutIntroMetricSvc,
		aboutIntroPointSvc,
		homepageServicesDocSvc,
		aboutIntroDocSvc,
	)

	// 6. Jalankan server
	addr := ":" + config.App.AppPort
	log.Printf("[server] Starting on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[server] Failed to start: %v", err)
	}
}
