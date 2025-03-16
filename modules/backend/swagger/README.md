# Swagger Documentation

This directory contains the Swagger documentation for the Software Engineer Salary Survey API.

## Generating Documentation

To generate or update the Swagger documentation, run:

```bash
make swagger
```

This will parse the API annotations in the code and generate the necessary Swagger files.

## Viewing Documentation

Once the application is running, you can access the Swagger UI at:

```text
http://localhost:4004/swagger/
```

## Adding Documentation to Endpoints

To add documentation to an endpoint, use the Swagger annotations in your handler functions. For example:

```go
// GetSurvey godoc
// @Summary Get a survey by ID
// @Description Get a survey by its ID
// @Tags surveys
// @Accept json
// @Produce json
// @Param id path string true "Survey ID"
// @Success 200 {object} SurveyResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/surveys/{id} [get]
func GetSurvey(c *fiber.Ctx) error {
    // Implementation
}
```

For more information on Swagger annotations, see the [Swaggo documentation](https://github.com/swaggo/swag#declarative-comments-format).
