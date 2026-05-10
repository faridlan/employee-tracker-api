package http

import (
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

// ==========================================
// HELPER: MAPPING ENTITY KE RESPONSE DTO
// ==========================================
func toAchievementResponse(ach *domain.Achievement) dto.AchievementResponse {
	if ach == nil {
		return dto.AchievementResponse{}
	}
	return dto.AchievementResponse{
		ID:          ach.ID,
		TargetID:    ach.TargetID,
		Nominal:     ach.Nominal,
		Description: ach.Description,
		ClosingDate: ach.ClosingDate,
		CreatedAt:   ach.CreatedAt,
		UpdatedAt:   ach.UpdatedAt,
	}
}

// RecordAchievement godoc
// @Summary Catat Pencapaian Target
// @Description Mencatat riwayat transaksi/pencairan (ledger) baru untuk sebuah target
// @Tags Achievement
// @Accept json
// @Produce json
// @Param request body dto.RecordAchievementRequest true "Payload data pencapaian"
// @Success 201 {object} utils.SuccessResponse[dto.AchievementResponse] "Pencapaian berhasil dicatat"
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

	input := domain.RecordAchievementInput{
		TargetID:    req.TargetID,
		Nominal:     req.Nominal,
		Description: req.Description,
		ClosingDate: req.ClosingDate,
	}

	result, err := h.usecase.RecordAchievement(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Pencapaian berhasil dicatat", toAchievementResponse(result))
}

// GetAchievementsByTarget godoc
// @Summary List Pencapaian Berdasarkan Target
// @Description Mengambil riwayat/ledger pencapaian untuk satu target spesifik
// @Tags Achievement
// @Produce json
// @Param target_id path string true "ID Target"
// @Success 200 {object} utils.SuccessResponse[[]dto.AchievementResponse] "Berhasil mengambil riwayat pencapaian"
// @Failure 400 {object} utils.ErrorResponse "Format UUID salah"
// @Failure 404 {object} utils.ErrorResponse "Target tidak ditemukan"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/targets/{target_id}/achievements [get]
func (h *AchievementHandler) GetAchievementsByTarget(c *fiber.Ctx) error {
	targetID := c.Params("target_id")
	if err := utils.ValidateUUID(targetID, "target_id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	results, err := h.usecase.GetAchievementsByTarget(c.Context(), targetID)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := make([]dto.AchievementResponse, 0)
	for _, ach := range results {
		res = append(res, toAchievementResponse(ach))
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil riwayat pencapaian", res)
}
