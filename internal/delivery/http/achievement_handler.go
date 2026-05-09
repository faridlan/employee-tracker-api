package http

import (
	// Sesuaikan path import dengan module-mu

	"github.com/faridlan/employee-tracker-api/internal/delivery/http/dto"
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type AchievementHandler struct {
	usecase domain.AchievementUsecase
}

func NewAchievementHandler(u domain.AchievementUsecase) *AchievementHandler {
	return &AchievementHandler{usecase: u}
}

// RecordAchievement godoc
// @Summary Catat Pencapaian Target
// @Description Mencatat riwayat transaksi/pencairan (ledger) baru untuk sebuah target
// @Tags Achievement
// @Accept json
// @Produce json
// @Param request body dto.RecordAchievementRequest true "Payload data pencapaian"
// @Success 201 {object} utils.SuccessResponse[utils.EmptyObj] "Pencapaian berhasil dicatat"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 404 {object} utils.ErrorResponse "Target tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/achievements [post]
func (h *AchievementHandler) RecordAchievement(c *fiber.Ctx) error {
	var req dto.RecordAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// Panggil Usecase dengan data dari DTO
	err := h.usecase.RecordAchievement(
		req.TargetID,
		req.Nominal,
		req.Description,
		req.ClosingDate,
	)

	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	// Untuk response ini, kita tidak perlu membalikkan data spesifik, cukup status sukses.
	// Jika ingin, kamu bisa membuat AchievementResponse DTO, tapi di aplikasi perbankan/ledger
	// biasanya HTTP 201 Created dengan pesan sukses sudah cukup.
	return utils.SendSuccess(c, fiber.StatusCreated, "Pencapaian berhasil dicatat", utils.EmptyObj{})
}
