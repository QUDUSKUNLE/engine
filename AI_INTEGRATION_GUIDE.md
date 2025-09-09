# AI Integration Guide for Diagnoxix

## Overview

This guide explains how AI has been integrated into your Diagnoxix healthcare application. The AI features provide intelligent analysis and interpretation capabilities to enhance patient care and diagnostic workflows.

## Features Implemented

### 1. Lab Result Interpretation
- **Endpoint**: `POST /v1/ai/interpret-lab`
- **Purpose**: Analyzes laboratory test results and provides medical interpretation
- **Input**: Lab test data with results and reference ranges
- **Output**: Comprehensive analysis including abnormaecommendations, and urgency level

### 2. Symptom Analysis
- **Endpoint**: `POST /v1/ai/analyze-symptoms`
- **Purpose**: Provides preliminary analysis of patient symptoms
- **Input**: List of symptoms, patient age, and gender
- **Output**: Possible conditions, urgency level, recommendations, and red flags

### 3. Medical Report Summarization
- **Endpoint**: `POST /v1/ai/summarize-report`
- **Purpose**: Creates patient-friendly or professional summaries of medical reports
- **Input**: Medical report text and summary type preference
- **Output**: Clear, organized summary in requested format

### 4. AI Capabilities Discovery
- **Endpoint**: `GET /v1/ai/capabilities`
- **Purpose**: Returns available AI features and their descriptions
- **Output**: List of available AI endpoints and capabilities

## API Usage Examples

### Lab Result Interpretation

```bash
curl -X POST http://localhost:7556/v1/ai/interpret-lab \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "patient_id": "patient-123",
    "diagnostic_centre_id": "center-456",
    "test_name": "Complete Blood Count (CBC)",
    "results": {
      "Hemoglobin": "8.5 g/dL",
      "White Blood Cells": "12,000 /μL",
      "Platelets": "150,000 /μL"
    },
    "reference_ranges": {
      "Hemoglobin": "12.0-15.5 g/dL",
      "White Blood Cells": "4,000-11,000 /μL",
      "Platelets": "150,000-450,000 /μL"
    }
  }'
```

### Symptom Analysis

```bash
curl -X POST http://localhost:7556/v1/ai/analyze-symptoms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "symptoms": ["persistent cough", "shortness of breath", "chest pain"],
    "age": 45,
    "gender": "male"
  }'
```

### Report Summarization

```bash
curl -X POST http://localhost:7556/v1/ai/summarize-report \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "medical_report": "CHEST X-RAY REPORT\n\nFindings: Bilateral lower lobe infiltrates...",
    "patient_friendly": true
  }'
```

## Configuration

### Environment Variables

Ensure your `.env` file contains:

```env
OPEN_API_KEY=your_openai_api_key_here
```

### Dependencies

The AI integration uses:
- OpenAI GPT-4o-mini model for cost-effective analysis
- Standard HTTP client with 30-second timeout
- JSON parsing for structured responses

## Architecture

### Service Layer
- `AIService` in `core/services/ai.services.go`
- Handles OpenAI API communication
- Provides structured response parsing
- Includes error handling and logging

### Handler Layer
- AI handlers in `adapters/handlers/ai.handler.go`
- Request validation and response formatting
- Integration with Echo framework
- Swagger documentation support

### Routes
- AI routes in `adapters/routes/ai.go`
- RESTful endpoint definitions
- Middleware integration

## Testing

Run the test script to verify AI integration:

```bash
go run test_ai_integration.go
```

This will test all three AI features with sample data.

## Security Considerations

1. **API Key Protection**: OpenAI API key is stored in environment variables
2. **JWT Authentication**: All AI endpoints require valid JWT tokens
3. **Input Validation**: All requests are validated before processing
4. **Rate Limiting**: Existing rate limiting middleware applies to AI endpoints
5. **Error Handling**: Sensitive information is not exposed in error messages

## Cost Management

- Uses GPT-4o-mini model for cost efficiency
- Temperature set to 0.3 for consistent medical analysis
- Token limits set to 1500 to control costs
- Timeout set to 30 seconds to prevent hanging requests

## Medical Disclaimers

All AI responses include appropriate medical disclaimers:
- "This analysis is for informational purposes only"
- "Should not replace professional medical consultation"
- "Always recommend consulting with healthcare providers"

## Integration with Existing AI System

Your application already has an extensive AI system in `adapters/ex/ai/`. The new AI service complements this by providing:

- Direct OpenAI integration for natural language processing
- Medical-specific prompt engineering
- Structured JSON responses for frontend integration
- RESTful API endpoints for easy consumption

## Future Enhancements

Consider these additional AI features:

1. **Drug Interaction Checking**: Analyze medication combinations
2. **Appointment Prioritization**: AI-driven scheduling based on urgency
3. **Medical Image Analysis**: Integration with your existing image analysis
4. **Treatment Recommendations**: Evidence-based treatment suggestions
5. **Risk Assessment**: Patient risk stratification
6. **Clinical Decision Support**: Integration with your existing decision support system

## Monitoring and Logging

- All AI requests are logged using your existing Zap logger
- Errors are captured with context for debugging
- Response times and token usage can be monitored
- Integration with your existing Prometheus metrics

## Support

For issues or questions about the AI integration:
1. Check the logs in `logs/medivue.log`
2. Verify OpenAI API key configuration
3. Test with the provided test script
4. Review API documentation at `/swagger/`
