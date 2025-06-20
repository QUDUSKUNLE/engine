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
	daysOfWeek := make([]string, len(dto.Slots))
	startTimes := make([]pgtype.Time, len(dto.Slots))
	endTimes := make([]pgtype.Time, len(dto.Slots))
	maxAppointments := make([]int32, len(dto.Slots))
	slotDurations := make([]int32, len(dto.Slots))
	breakTimes := make([]int32, len(dto.Slots))

	for i, slot := range dto.Slots {
		diagnosticCentreIDs[i] = dto.DiagnosticCentreID
		daysOfWeek[i] = slot.DayOfWeek

		// Convert time.Time to pgtype.Time for start time
		startTime := pgtype.Time{}
		// Extract only the time part (hours, minutes, seconds, microseconds)
		hour, min, sec := slot.StartTime.Clock()
		nsec := slot.StartTime.Nanosecond()
		startTime.Microseconds = int64(hour*3600+min*60+sec)*1000000 + int64(nsec/1000)
		startTime.Valid = true
		startTimes[i] = startTime

		// Convert time.Time to pgtype.Time for end time
		endTime := pgtype.Time{}
		// Extract only the time part (hours, minutes, seconds, microseconds)
		hour, min, sec = slot.EndTime.Clock()
		nsec = slot.EndTime.Nanosecond()
		endTime.Microseconds = int64(hour*3600+min*60+sec)*1000000 + int64(nsec/1000)
		endTime.Valid = true
		endTimes[i] = endTime

		maxAppointments[i] = slot.MaxAppointments
		slotDurations[i] = slot.SlotDuration // Already an int32
		breakTimes[i] = slot.BreakTime       // Already an int32
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
	_, err := PrivateMiddlewareContext(ctx, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return nil, err
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
	_, err := PrivateMiddlewareContext(ctx, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return nil, err
	}

	diagnosticCentreID := ctx.Param("diagnostic_centre_id")
	dayOfWeek := ctx.Param("day_of_week")

	dto, _ := ctx.Get(utils.ValidatedBodyDTO).(*domain.UpdateAvailabilityDTO)
	if dto == nil {
		return nil, fmt.Errorf("invalid request body")
	}
	startTime := pgtype.Time{}
	if dto.StartTime != nil {
		if err := startTime.Scan(*dto.StartTime); err != nil {
			return nil, fmt.Errorf("%s: %w", errInvalidEndTime, err)
		}
	}

	endTime := pgtype.Time{}
	if dto.EndTime != nil {
		if err := endTime.Scan(*dto.EndTime); err != nil {
			return nil, fmt.Errorf("%s: %w", errInvalidEndTime, err)
		}
	}

	var slotDuration int32
	if dto.SlotDuration != nil {
		slotDuration = *dto.SlotDuration
	}

	var maxAppointments int32
	if dto.MaxAppointments != nil {
		maxAppointments = *dto.MaxAppointments
	}

	var breakTime int32
	if dto.BreakTime != nil {
		breakTime = *dto.BreakTime
	}

	return &db.Update_AvailabilityParams{
		DiagnosticCentreID: diagnosticCentreID,
		SlotDuration:       slotDuration,
		MaxAppointments:    pgtype.Int4{Int32: maxAppointments, Valid: true},
		BreakTime:          pgtype.Int4{Int32: breakTime},
		StartTime:          startTime,
		EndTime:            endTime,
		DayOfWeek:          dayOfWeek,
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
	_, err := PrivateMiddlewareContext(ctx, []db.UserEnum{db.UserEnumDIAGNOSTICCENTREOWNER, db.UserEnumDIAGNOSTICCENTREMANAGER})
	if err != nil {
		return nil, err
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
	daysOfWeek := make([]string, len(dto.Slots))
	startTimes := make([]pgtype.Time, len(dto.Slots))
	endTimes := make([]pgtype.Time, len(dto.Slots))
	maxAppointments := make([]int32, len(dto.Slots))
	slotDurations := make([]int32, len(dto.Slots))
	breakTimes := make([]int32, len(dto.Slots))

	for i, slot := range dto.Slots {
		diagnosticCentreIDs[i] = slot.DiagnosticCentreID
		daysOfWeek[i] = string(slot.DayOfWeek)

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
			slotDurations[i] = *slot.SlotDuration // Already an int32
		}

		if slot.BreakTime != nil {
			breakTimes[i] = *slot.BreakTime // Already an int32
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
	diagnosticCentreID := ctx.Param("diagnostic_centre_id")
	if diagnosticCentreID == "" {
		return nil, fmt.Errorf("diagnostic_centre_id is required")
	}
	// Day of week is optional in query params
	dayOfWeek := ctx.QueryParam("day_of_week")
	if dayOfWeek != "" {
		// Validate day_of_week if provided
		validDays := map[string]bool{
			"monday": true, "tuesday": true, "wednesday": true,
			"thursday": true, "friday": true, "saturday": true, "sunday": true,
		}
		if !validDays[dayOfWeek] {
			return nil, fmt.Errorf("invalid day of week")
		}
	}
	return &db.Get_AvailabilityParams{
		DiagnosticCentreID: diagnosticCentreID,
		Column2:            dayOfWeek,
	}, nil
}

// GetAvailability retrieves availability slots for a diagnostic centre
func (s *ServicesHandler) GetAvailability(ctx echo.Context) error {
	params, err := s.validateGetAvailabilityInput(ctx)
	if err != nil {
		return utils.ErrorResponse(http.StatusBadRequest, err, ctx)
	}
	var slots []*db.DiagnosticCentreAvailability
	if params.Column2 == "" {
		slots, err = s.AvailabilityRepo.GetDiagnosticAvailability(ctx.Request().Context(), params.DiagnosticCentreID)
		if err != nil {
			return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
		}
	} else {
		slots, err = s.AvailabilityRepo.GetAvailability(ctx.Request().Context(), *params)
		if err != nil {
			return utils.ErrorResponse(http.StatusInternalServerError, err, ctx)
		}
	}
	if len(slots) == 0 {
		return utils.ResponseMessage(http.StatusOK, []struct{}{}, ctx)
	}

	result := make([]*domain.AvailabilitySlot, len(slots))
	for i, slot := range slots {
		result[i] = &domain.AvailabilitySlot{
			ID:                 slot.ID,
			DiagnosticCentreID: slot.DiagnosticCentreID,
			DayOfWeek:          slot.DayOfWeek,
			StartTime:          fmt.Sprintf("%02d:%02d:%02d", slot.StartTime.Microseconds/3600000000, (slot.StartTime.Microseconds%3600000000)/60000000, ((slot.StartTime.Microseconds%3600000000)%60000000)/1000000),
			EndTime:            fmt.Sprintf("%02d:%02d:%02d", slot.EndTime.Microseconds/3600000000, (slot.EndTime.Microseconds%3600000000)/60000000, ((slot.EndTime.Microseconds%3600000000)%60000000)/1000000),
			SlotDuration:       fmt.Sprintf("%d mins", slot.SlotDuration),
			CreatedAt:          slot.CreatedAt.Time.String(),
			UpdatedAt:          slot.UpdatedAt.Time.String(),
			BreakTime:          fmt.Sprintf("%d mins", slot.BreakTime.Int32),
			MaxAppointments:    slot.MaxAppointments.Int32,
		}
	}

	return utils.ResponseMessage(http.StatusOK, result, ctx)
}
