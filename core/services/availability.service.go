package services

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/medicue/adapters/db"
	"github.com/medicue/core/domain"
	"github.com/medicue/core/utils"
)

const (
	errInvalidStartTime    = "invalid start time format"
	errInvalidEndTime      = "invalid end time format"
	errInvalidSlotDuration = "invalid slot duration format"
	errInvalidBreakTime    = "invalid break time format"
)

// convertSlotToArrays converts a slice of availability slots to arrays of individual fields
func (s *ServicesHandler) convertSlotToArrays(dto *domain.CreateAvailabilityDTO) (*db.Create_AvailabilityParams, error) {
	diagnosticCentreIDs := make([]string, len(dto.Slots))
	daysOfWeek := make([]db.Weekday, len(dto.Slots))
	startTimes := make([]pgtype.Time, len(dto.Slots))
	endTimes := make([]pgtype.Time, len(dto.Slots))
	maxAppointments := make([]int32, len(dto.Slots))
	slotDurations := make([]pgtype.Interval, len(dto.Slots))
	breakTimes := make([]pgtype.Interval, len(dto.Slots))

	for i, slot := range dto.Slots {
		diagnosticCentreIDs[i] = dto.DiagnosticCentreID
		daysOfWeek[i] = slot.DayOfWeek

		startTime := pgtype.Time{}
		if err := startTime.Scan(slot.StartTime); err != nil {
			return nil, fmt.Errorf("%s: %w", errInvalidStartTime, err)
		}
		startTimes[i] = startTime

		endTime := pgtype.Time{}
		if err := endTime.Scan(slot.EndTime); err != nil {
			return nil, fmt.Errorf("%s: %w", errInvalidEndTime, err)
		}
		endTimes[i] = endTime

		maxAppointments[i] = int32(slot.MaxAppointments)

		slotDuration := pgtype.Interval{}
		if err := slotDuration.Scan(slot.SlotDuration); err != nil {
			return nil, fmt.Errorf("%s: %w", errInvalidSlotDuration, err)
		}
		slotDurations[i] = slotDuration

		breakTime := pgtype.Interval{}
		if err := breakTime.Scan(slot.BreakTime); err != nil {
			return nil, fmt.Errorf("%s: %w", errInvalidBreakTime, err)
		}
		breakTimes[i] = breakTime
	}

	return &db.Create_AvailabilityParams{
		Column1: diagnosticCentreIDs,
		Column2: daysOfWeek,
		Column3: startTimes,
		Column4: endTimes,
		Column5: maxAppointments,
		Column6: slotDurations,
		Column7: breakTimes,
	}, nil
}

// validateCreateAvailabilityInput validates the input for creating availability slots
func (s *ServicesHandler) validateCreateAvailabilityInput(ctx echo.Context) (*domain.CreateAvailabilityDTO, error) {
	// Authenticate and authorize user - owner or manager only
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return nil, err
		}
	}

	// This validated at the middleware level
	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.CreateAvailabilityDTO)
	if dto == nil {
		return nil, fmt.Errorf("invalid request body")
	}

	return dto, nil
}

// CreateAvailability creates new availability slots for a diagnostic centre
func (s *ServicesHandler) CreateAvailability(ctx echo.Context) error {
	dto, err := s.validateCreateAvailabilityInput(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
	}

	params, err := s.convertSlotToArrays(dto)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	slot, err := s.AvailabilityRepo.CreateAvailability(ctx.Request().Context(), *params)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusCreated, slot)
}

// validateUpdateAvailabilityInput validates and converts the input for updating availability slots
func (s *ServicesHandler) validateUpdateAvailabilityInput(ctx echo.Context) (*db.Update_AvailabilityParams, error) {
	// Authenticate and authorize user
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return nil, err
		}
	}

	diagnosticCentreID := ctx.Param("diagnostic_centre_id")
	dayOfWeek := ctx.Param("day_of_week")

	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.UpdateAvailabilityDTO)
	if dto == nil {
		return nil, fmt.Errorf("invalid request body")
	}

	// Convert slot duration to pgtype.Interval
	slotDuration := pgtype.Interval{}
	if err := slotDuration.Scan(dto.SlotDuration); err != nil {
		return nil, fmt.Errorf("%s: %w", errInvalidSlotDuration, err)
	}

	// Convert break time to pgtype.Interval
	breakTime := pgtype.Interval{}
	if err := breakTime.Scan(dto.BreakTime); err != nil {
		return nil, fmt.Errorf("%s: %w", errInvalidBreakTime, err)
	}

	// Convert end time to pgtype.Time
	endTime := pgtype.Time{}
	if err := endTime.Scan(*dto.EndTime); err != nil {
		return nil, fmt.Errorf("%s: %w", errInvalidEndTime, err)
	}

	return &db.Update_AvailabilityParams{
		DiagnosticCentreID: diagnosticCentreID,
		SlotDuration:       slotDuration,
		MaxAppointments:    pgtype.Int4{Int32: int32(*dto.MaxAppointments), Valid: true},
		BreakTime:          breakTime,
		EndTime:            endTime,
		DayOfWeek:          db.Weekday(dayOfWeek),
	}, nil
}

// UpdateAvailability updates an existing availability slot
func (s *ServicesHandler) UpdateAvailability(ctx echo.Context) error {
	params, err := s.validateUpdateAvailabilityInput(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
	}

	slot, err := s.AvailabilityRepo.UpdateAvailability(ctx.Request().Context(), *params)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusOK, slot)
}

// validateUpdateManyAvailabilityInput validates the input for updating multiple availability slots
func (s *ServicesHandler) validateUpdateManyAvailabilityInput(ctx echo.Context) (*domain.UpdateManyAvailabilityDTO, error) {
	// Authenticate and authorize user
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return nil, err
		}
	}

	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.UpdateManyAvailabilityDTO)
	if dto == nil {
		return nil, fmt.Errorf("invalid request body")
	}

	return dto, nil
}

// convertUpdateManySlotToArrays converts a slice of update availability slots to arrays of individual fields
func (s *ServicesHandler) convertUpdateManySlotToArrays(dto *domain.UpdateManyAvailabilityDTO) (db.Update_Many_AvailabilityParams, error) {
	diagnosticCentreIDs := make([]string, len(dto.Slots))
	daysOfWeek := make([]db.Weekday, len(dto.Slots))
	startTimes := make([]pgtype.Time, len(dto.Slots))
	endTimes := make([]pgtype.Time, len(dto.Slots))
	maxAppointments := make([]int32, len(dto.Slots))
	slotDurations := make([]pgtype.Interval, len(dto.Slots))
	breakTimes := make([]pgtype.Interval, len(dto.Slots))

	for i, slot := range dto.Slots {
		diagnosticCentreIDs[i] = slot.DiagnosticCentreID
		daysOfWeek[i] = slot.DayOfWeek

		if slot.StartTime != nil {
			startTime := pgtype.Time{}
			if err := startTime.Scan(*slot.StartTime); err != nil {
				return db.Update_Many_AvailabilityParams{}, fmt.Errorf("%s: %w", errInvalidStartTime, err)
			}
			startTimes[i] = startTime
		}

		if slot.EndTime != nil {
			endTime := pgtype.Time{}
			if err := endTime.Scan(*slot.EndTime); err != nil {
				return db.Update_Many_AvailabilityParams{}, fmt.Errorf("%s: %w", errInvalidEndTime, err)
			}
			endTimes[i] = endTime
		}

		if slot.MaxAppointments != nil {
			maxAppointments[i] = int32(*slot.MaxAppointments)
		}

		if slot.SlotDuration != nil {
			slotDuration := pgtype.Interval{}
			if err := slotDuration.Scan(*slot.SlotDuration); err != nil {
				return db.Update_Many_AvailabilityParams{}, fmt.Errorf("%s: %w", errInvalidSlotDuration, err)
			}
			slotDurations[i] = slotDuration
		}

		if slot.BreakTime != nil {
			breakTime := pgtype.Interval{}
			if err := breakTime.Scan(*slot.BreakTime); err != nil {
				return db.Update_Many_AvailabilityParams{}, fmt.Errorf("%s: %w", errInvalidBreakTime, err)
			}
			breakTimes[i] = breakTime
		}
	}

	return db.Update_Many_AvailabilityParams{
		Column1: diagnosticCentreIDs,
		Column2: daysOfWeek,
		Column3: startTimes,
		Column4: endTimes,
		Column5: maxAppointments,
		Column6: slotDurations,
		Column7: breakTimes,
	}, nil
}

// UpdateManyAvailability updates multiple availability slots in bulk
func (s *ServicesHandler) UpdateManyAvailability(ctx echo.Context) error {
	dto, err := s.validateUpdateManyAvailabilityInput(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
	}

	params, err := s.convertUpdateManySlotToArrays(dto)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	slots, err := s.AvailabilityRepo.UpdateManyAvailability(ctx.Request().Context(), params)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusOK, slots)
}

// validateGetAvailabilityInput validates the input for getting availability slots
func (s *ServicesHandler) validateGetAvailabilityInput(ctx echo.Context) (*db.Get_AvailabilityParams, error) {
	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.GetAvailabilityDTO)
	if dto == nil {
		return nil, fmt.Errorf("invalid request body")
	}

	return &db.Get_AvailabilityParams{
		DiagnosticCentreID: dto.DiagnosticCentreID,
		Column2:            db.Weekday(dto.DayOfWeek),
	}, nil
}

// GetAvailability retrieves availability slots for a diagnostic centre
func (s *ServicesHandler) GetAvailability(ctx echo.Context) error {
	params, err := s.validateGetAvailabilityInput(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}

	slots, err := s.AvailabilityRepo.GetAvailability(ctx.Request().Context(), *params)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.JSON(http.StatusOK, slots)
}

// validateDeleteAvailabilityInput validates the input for deleting availability slots
func (s *ServicesHandler) validateDeleteAvailabilityInput(ctx echo.Context) (*db.Delete_AvailabilityParams, error) {
	// Authenticate and authorize user
	_, err := utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREOWNER))
	if err != nil {
		_, err = utils.PrivateMiddlewareContext(ctx, string(db.UserEnumDIAGNOSTICCENTREMANAGER))
		if err != nil {
			return nil, err
		}
	}

	return &db.Delete_AvailabilityParams{
		DiagnosticCentreID: ctx.Param("diagnostic_centre_id"),
		DayOfWeek:          db.Weekday(ctx.Param("day_of_week")),
	}, nil
}

// DeleteAvailability removes an availability slot for a diagnostic centre
func (s *ServicesHandler) DeleteAvailability(ctx echo.Context) error {
	params, err := s.validateDeleteAvailabilityInput(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusUnauthorized, err, ctx)
	}

	err = s.AvailabilityRepo.DeleteAvailability(ctx.Request().Context(), *params)
	if err != nil {
		return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}
