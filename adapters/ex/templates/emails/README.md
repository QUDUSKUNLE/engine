# Email Templates

This package provides a collection of responsive, modern email templates for the DiagnoxixAI healthcare platform.

## Available Templates

- Appointment Confirmation
- Appointment Cancellation
- Appointment Reschedule
- Appointment Reminder
- Payment Confirmation
- Test Results Available
- Staff Notification
- Policy Update

## Usage

```go
// Initialize the template handler
handler := emails.NewEmailTemplateHandler()

// Create template data
data := &emails.AppointmentData{
    PatientName:     "John Doe",
    AppointmentID:   "APT123",
    AppointmentDate: time.Now(),
    TimeSlot:        "10:00 AM",
    CentreName:      "Medical Centre",
    TestType:        "BLOOD_TEST",
    Notes:           "Please fast for 8 hours before the test",
}

// Render the template
html, err := handler.RenderAppointmentConfirmation(data)
if err != nil {
    log.Fatal(err)
}

// Use the rendered HTML in your email service
// ...
```

## Template Features

- Responsive design
- Modern styling with gradients and shadows
- Consistent branding across all templates
- Support for emojis as icons
- Clear information hierarchy
- Properly formatted dates and currency
- Mobile-friendly layout

## Data Structures

Each template has its own data structure that defines the required fields:

- `AppointmentData`: For appointment-related emails
- `PaymentData`: For payment confirmation emails
- `TestResultsData`: For test results notifications
- `StaffNotificationData`: For staff notifications
- `PolicyUpdateData`: For policy update notifications

## Customization

The templates use CSS variables for easy customization:

```css
:root {
    --primary-color: #2563eb;
    --secondary-color: #1e40af;
    --accent-color: #60a5fa;
    --text-primary: #1f2937;
    --text-secondary: #4b5563;
    --background: #f3f4f6;
}
```

## Testing

Run the tests using:

```bash
go test ./adapters/ex/templates/emails -v
```
