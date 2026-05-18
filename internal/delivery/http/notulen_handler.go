package http

import (
	"strings"

	"github.com/faridlan/employee-tracker-api/internal/delivery/http/dto"
	"github.com/faridlan/employee-tracker-api/internal/domain"
	"github.com/faridlan/employee-tracker-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type MeetingMinuteHandler struct {
	usecase domain.MeetingMinuteUsecase
}

func NewMeetingMinuteHandler(u domain.MeetingMinuteUsecase) *MeetingMinuteHandler {
	return &MeetingMinuteHandler{usecase: u}
}

// ==========================================
// HELPER: MAPPING ENTITY KE RESPONSE DTO
// ==========================================
func toMeetingMinuteResponse(m *domain.MeetingMinute) dto.MeetingMinuteResponse {
	if m == nil {
		return dto.MeetingMinuteResponse{}
	}

	// PROSES PECAH STRING (DARI DB) MENJADI ARRAY UNTUK FRONTEND
	var extParticipants []string
	if m.ExternalParticipants != nil && *m.ExternalParticipants != "" {
		// Cara baru: Langsung looping menggunakan SplitSeq tanpa array perantara
		for name := range strings.SplitSeq(*m.ExternalParticipants, ",") {
			extParticipants = append(extParticipants, strings.TrimSpace(name))
		}
	} else {
		extParticipants = []string{}
	}

	res := dto.MeetingMinuteResponse{
		ID:                   m.ID,
		Division:             m.Division,
		Title:                m.Title,
		MeetingDate:          m.MeetingDate,
		MeetingType:          m.MeetingType,
		Summary:              m.Summary,
		Notes:                m.Notes,
		Speaker:              m.Speaker,
		NumberOfParticipants: m.NumberOfParticipants,
		ExternalParticipants: extParticipants, // DIMASUKKAN KE SINI
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}

	// Map Participants
	for _, p := range m.Participants {
		partRes := dto.MeetingParticipantResponse{
			ID:         p.ID,
			MinuteID:   p.MinuteID,
			EmployeeID: p.EmployeeID,
			CreatedAt:  p.CreatedAt,
		}
		if p.Employee != nil {
			emp := toEmployeeResponse(p.Employee)
			partRes.Employee = &emp
		}
		res.Participants = append(res.Participants, partRes)
	}

	// Map Results / Tasks
	for _, r := range m.Results {
		resRes := dto.MeetingResultResponse{
			ID:                   r.ID,
			MinuteID:             r.MinuteID,
			EmployeeID:           r.EmployeeID,
			TargetDescription:    r.TargetDescription,
			TargetNominal:        r.TargetNominal,
			AchievementStatus:    r.AchievementStatus,
			TargetCompletionDate: r.TargetCompletionDate,
			CreatedAt:            r.CreatedAt,
			UpdatedAt:            r.UpdatedAt,
		}
		if r.Employee != nil {
			emp := toEmployeeResponse(r.Employee)
			resRes.Employee = &emp
		}
		res.Results = append(res.Results, resRes)
	}

	// Map Images
	for _, i := range m.Images {
		res.Images = append(res.Images, dto.MeetingImageResponse{
			ID:        i.ID,
			MinuteID:  i.MinuteID,
			FileURL:   i.FileURL,
			CreatedAt: i.CreatedAt,
		})
	}

	return res
}

// CreateMeeting godoc
// @Summary Buat Notulen Rapat Baru
// @Description Menyimpan notulen rapat beserta daftar peserta, tugas (action items), dan dokumentasi
// @Tags Meeting
// @Accept json
// @Produce json
// @Param request body dto.CreateMeetingRequest true "Payload data notulen"
// @Success 201 {object} utils.SuccessResponse[dto.MeetingMinuteResponse] "Notulen berhasil dibuat"
// @Failure 400 {object} utils.ErrorResponse "Format JSON salah atau validasi gagal"
// @Failure 500 {object} utils.ErrorResponse "Gagal menyimpan data notulen"
// @Router /api/meetings [post]
func (h *MeetingMinuteHandler) CreateMeeting(c *fiber.Ctx) error {
	var req dto.CreateMeetingRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// PROSES GABUNG ARRAY (DARI FE) MENJADI SATU STRING UNTUK DB
	var externalNames *string
	if len(req.ExternalParticipants) > 0 {
		joinedNames := strings.Join(req.ExternalParticipants, ", ")
		externalNames = &joinedNames
	}

	// Mapping DTO ke Domain Input
	input := domain.CreateMeetingInput{
		Division:             req.Division,
		Title:                req.Title,
		MeetingDate:          req.MeetingDate,
		MeetingType:          req.MeetingType,
		Summary:              req.Summary,
		Notes:                req.Notes,
		Speaker:              req.Speaker,
		ExternalParticipants: externalNames, // DIMASUKKAN KE SINI
		ParticipantIDs:       req.ParticipantIDs,
		ImageURLs:            req.ImageURLs,
	}

	for _, r := range req.Results {
		input.Results = append(input.Results, domain.CreateMeetingResultInput{
			EmployeeID:           r.EmployeeID,
			TargetDescription:    r.TargetDescription,
			TargetNominal:        r.TargetNominal,
			TargetCompletionDate: r.TargetCompletionDate,
		})
	}

	result, err := h.usecase.CreateMeeting(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	res := toMeetingMinuteResponse(result)
	return utils.SendSuccess(c, fiber.StatusCreated, "Notulen berhasil dibuat", res)
}

// UpdateMeeting godoc
// @Summary Update Data Rapat
// @Description Memperbaharui informasi dasar dari notulen rapat
// @Tags Meeting
// @Accept json
// @Produce json
// @Param id path string true "ID Notulen"
// @Param request body dto.UpdateMeetingRequest true "Payload data notulen"
// @Success 200 {object} utils.SuccessResponse[dto.MeetingMinuteResponse] "Notulen berhasil diperbaharui"
// @Router /api/meetings/{id} [put]
func (h *MeetingMinuteHandler) UpdateMeeting(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.UpdateMeetingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	// PROSES GABUNG ARRAY (DARI FE) MENJADI SATU STRING UNTUK DB
	var externalNames *string
	if len(req.ExternalParticipants) > 0 {
		joinedNames := strings.Join(req.ExternalParticipants, ", ")
		externalNames = &joinedNames
	}

	input := domain.UpdateMeetingInput{
		ID:                   id,
		Division:             req.Division,
		Title:                req.Title,
		MeetingDate:          req.MeetingDate,
		MeetingType:          req.MeetingType,
		Summary:              req.Summary,
		Notes:                req.Notes,
		Speaker:              req.Speaker,
		ExternalParticipants: externalNames, // DIMASUKKAN KE SINI
	}

	result, err := h.usecase.UpdateMeeting(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Notulen berhasil diperbaharui", toMeetingMinuteResponse(result))
}

// GetMeetingDetail godoc
// @Summary Detail Notulen Rapat
// @Description Menampilkan detail notulen beserta peserta dan tugas-tugasnya
// @Tags Meeting
// @Produce json
// @Param id path string true "ID Notulen"
// @Success 200 {object} utils.SuccessResponse[dto.MeetingMinuteResponse] "Berhasil mengambil detail notulen"
// @Router /api/meetings/{id} [get]
func (h *MeetingMinuteHandler) GetMeetingDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.ValidateUUID(id, "id"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.usecase.GetMeetingDetails(c.Context(), id)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil detail notulen", toMeetingMinuteResponse(result))
}

// GetAllMeetings godoc
// @Summary Data Semua Notulen
// @Description Menampilkan semua data notulen rapat terbaru
// @Tags Meeting
// @Produce json
// @Success 200 {object} utils.SuccessResponse[[]dto.MeetingMinuteResponse] "Berhasil mengambil data notulen"
// @Router /api/meetings [get]
func (h *MeetingMinuteHandler) GetAllMeetings(c *fiber.Ctx) error {
	results, err := h.usecase.GetAllMeetings(c.Context())
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	var res []dto.MeetingMinuteResponse
	for _, m := range results {
		res = append(res, toMeetingMinuteResponse(m))
	}

	// Fiber akan mereturn `null` jika array kosong. Agar aman di FE, kembalikan slice kosong jika length 0
	if len(res) == 0 {
		res = []dto.MeetingMinuteResponse{}
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Berhasil mengambil data notulen", res)
}

// UpdateTaskStatus godoc
// @Summary Update Status Tugas (Action Item)
// @Description Endpoint khusus untuk karyawan mengupdate status tugas rapat mereka
// @Tags Meeting
// @Accept json
// @Produce json
// @Param resultId path string true "ID Result / Tugas"
// @Param request body dto.UpdateResultStatusRequest true "Payload status"
// @Success 200 {object} utils.SuccessResponse[dto.MeetingResultResponse] "Status tugas berhasil diperbaharui"
// @Router /api/meetings/results/{resultId}/status [patch]
func (h *MeetingMinuteHandler) UpdateTaskStatus(c *fiber.Ctx) error {
	resultID := c.Params("resultId")
	if err := utils.ValidateUUID(resultID, "resultId"); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req dto.UpdateResultStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format JSON salah")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	input := domain.UpdateResultStatusInput{
		ResultID:          resultID,
		AchievementStatus: req.AchievementStatus,
	}

	result, err := h.usecase.UpdateTaskStatus(c.Context(), input)
	if err != nil {
		return utils.HandleDomainError(c, err)
	}

	// Convert langsung result ke DTO
	res := dto.MeetingResultResponse{
		ID:                   result.ID,
		MinuteID:             result.MinuteID,
		EmployeeID:           result.EmployeeID,
		TargetDescription:    result.TargetDescription,
		TargetNominal:        result.TargetNominal,
		AchievementStatus:    result.AchievementStatus,
		TargetCompletionDate: result.TargetCompletionDate,
		UpdatedAt:            result.UpdatedAt,
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Status tugas berhasil diperbaharui", res)
}
