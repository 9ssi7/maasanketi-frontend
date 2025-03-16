# Survey CLI Tool

This command-line tool allows developers to create surveys in the Software Engineer Salary Survey platform.

## Usage

### Using Default Template

```bash
./surveycli --default --title "Your Survey Title" --description "Your survey description" --min-time 5
```

Optional parameters:

- `--slug "custom-slug"`: Specify a custom slug for the survey URL (if not provided, one will be generated from the title)

### Using Custom JSON File

```bash
./surveycli --file path/to/your/survey.json
```

Optional parameters:

- `--slug "custom-slug"`: Specify a custom slug for the survey URL (if not provided, one will be generated from the title)

## JSON Format

The JSON file should follow this format:

```json
{
  "title": "Your Survey Title",
  "description": "Your survey description",
  "minCompletionTimeMin": 5,
  "questions": [
    {
      "id": "question_id",
      "text": "Question text",
      "type": "select",
      "required": true,
      "order": 1,
      "options": [
        {"id": "option1", "value": "Option 1"},
        {"id": "option2", "value": "Option 2"}
      ]
    },
    // More questions...
  ]
}
```

### Question Types

- `text`: Free text input
- `number`: Numeric input
- `select`: Single select dropdown
- `multi-select`: Multiple select checkboxes
- `boolean`: Yes/No question

### Conditional Questions

You can make questions conditional on previous answers:

```json
{
  "id": "conditional_question",
  "text": "This question only appears if condition is met",
  "type": "text",
  "required": true,
  "order": 3,
  "conditional": {
    "questionId": "previous_question_id",
    "value": "option1,option2"
  }
}
```

## Sample JSON

A sample JSON file is provided in `sample_survey.json` that you can use as a template for creating custom surveys.

## Environment Variables

The CLI tool requires the following environment variables:

- `DB_CONN_STR`: PostgreSQL connection string

Example:

```bash
export DB_CONN_STR="postgres://username:password@localhost:5432/database_name"
```

## Features

### Graceful Shutdown

The CLI tool supports graceful shutdown when receiving SIGINT (Ctrl+C) or SIGTERM signals. When a shutdown signal is received, the tool will:

1. Log a message indicating it's shutting down
2. Cancel any ongoing operations
3. Wait for 2 seconds to allow operations to complete
4. Exit cleanly

### Database Operation Timeout

Database operations have a 30-second timeout to prevent the tool from hanging indefinitely if there are connection issues. If a database operation takes longer than 30 seconds, the tool will exit with an error message.
