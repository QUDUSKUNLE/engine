package services

import (
	"net/http"

		"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
	"github.com/medicue/adapters/db"
	"github.com/labstack/echo/v4"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

// CreateAvailability creates a new availability slot for a diagnostic centre
func (s *ServicesHandler) CreateAvailability(ctx echo.Context) error {
	var req domain.CreateAvailabilityDTO

	// Authenticate and authorize user - owner or manager only
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
		}
	}
		// This validated at the middleware level
	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.CreateAvailabilityDTO)

	// Parse diagnostic centre ID
	centreID, err := uuid.Parse(req.DiagnosticCentreID)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	// Create availability slot
	startTime := pgtype.Time{}
	if err := startTime.Scan(dto.StartTime); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	endTime := pgtype.Time{}
	if err := endTime.Scan(dto.EndTime); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	slotDuration := pgtype.Interval{}
	if err := slotDuration.Scan(dto.SlotDuration); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	breakTime := pgtype.Interval{}
	if err := breakTime.Scan(dto.BreakTime); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	slot, err := s.AvailabilityRepo.CreateAvailability(ctx.Request().Context(), db.Create_AvailabilityParams{
		DiagnosticCentreID: centreID.String(),
		DayOfWeek:          dto.DayOfWeek,
		StartTime:          startTime,
		EndTime:            endTime,
		MaxAppointments:    pgtype.Int4{Int32: int32(dto.MaxAppointments), Valid: true},
		SlotDuration:       slotDuration,
		BreakTime:          breakTime,
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusCreated, slot)
}

// UpdateAvailability updates an existing availability slot
func (s *ServicesHandler) UpdateAvailability(ctx echo.Context) error {
	// Authenticate and authorize user - owner or manager only
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
		}
	}

	diagnosticCentreID := ctx.Param("diagnostic_centre_id")
	dayOfWeek := ctx.Param("day_of_week")

		// This validated at the middleware level
	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.UpdateAvailabilityDTO)

	// Convert slot duration to pgtype.Interval
	slotDuration := pgtype.Interval{}
	if err := slotDuration.Scan(dto.SlotDuration); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	// Convert break time to pgtype.Interval
	breakTime := pgtype.Interval{}
	if err := breakTime.Scan(dto.BreakTime); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	// Convert end time to pgtype.Time
	endTime := pgtype.Time{}
	if err := endTime.Scan(*dto.EndTime); err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	// Update availability
	slot, err := s.AvailabilityRepo.UpdateAvailability(ctx.Request().Context(), db.Update_AvailabilityParams{
		DiagnosticCentreID: diagnosticCentreID,
		SlotDuration: slotDuration,
		MaxAppointments: pgtype.Int4{Int32: int32(*dto.MaxAppointments), Valid: true},
		BreakTime: breakTime,
		EndTime: endTime,
		DayOfWeek: db.Weekday(dayOfWeek),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusOK, slot)
}

// GetAvailability retrieves availability slots for a diagnostic centre
func (s *ServicesHandler) GetAvailability(ctx echo.Context) error {

			// This validated at the middleware level
	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.GetAvailabilityDTO)

	// Get availability slots
	slots, err := s.AvailabilityRepo.GetAvailability(ctx.Request().Context(), db.Get_AvailabilityParams{
		DiagnosticCentreID: dto.DiagnosticCentreID,
		Column2: db.Weekday(dto.DayOfWeek),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusOK, slots)
}

// DeleteAvailability removes an availability slot for a diagnostic centre
func (s *ServicesHandler) DeleteAvailability(ctx echo.Context) error {
	// Authenticate and authorize user - owner or manager only
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
		}
	}

	diagnosticCentreID := ctx.Param("diagnostic_centre_id")
	dayOfWeek := ctx.Param("day_of_week")

	// Delete availability
	err = s.AvailabilityRepo.DeleteAvailability(ctx.Request().Context(), db.Delete_AvailabilityParams{
		DiagnosticCentreID: diagnosticCentreID,
		DayOfWeek: db.Weekday(dayOfWeek),
	})
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}
