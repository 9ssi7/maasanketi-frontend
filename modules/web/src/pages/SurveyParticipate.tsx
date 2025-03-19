import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Formik, Form, Field, ErrorMessage, FormikHelpers } from "formik";
import * as Yup from "yup";
import "../styles/SurveyParticipate.css";
import {
  surveyGet,
  surveyResponseComplete,
  surveyResponseStart,
} from "../services/survey.api";
import {
  SurveyAnswers,
  Survey,
  SurveyResponseStartResponse,
} from "../types/survey.types";
import { isSuccess } from "../services/base.api";

const SurveyParticipate = () => {
  const { slug } = useParams<{ slug: string }>();
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [step, setStep] = useState(1);
  const [isAnonymous, setIsAnonymous] = useState(true);
  const [responseId, setResponseId] = useState<string | null>("");

  useEffect(() => {
    if (!slug) return;
    surveyGet(slug)
      .then(([survey, status]) => {
        if (isSuccess(status)) {
          setSurvey(survey as Survey);
        } else {
          setError("Failed to load survey. Please try again later.");
        }
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [slug]);

  const startSurvey = async () => {
    if (!slug) return;

    try {
      setIsLoading(true);
      const [response, status] = await surveyResponseStart(slug, isAnonymous);

      if (isSuccess(status) && response) {
        setResponseId((response as SurveyResponseStartResponse).id);
        setStep(2);
      } else {
        setError("Failed to start survey. Please try again later.");
      }
    } catch (err) {
      setError("Failed to start survey. Please try again later.");
      console.error("Error starting survey:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (
    values: SurveyAnswers,
    { setSubmitting }: FormikHelpers<SurveyAnswers>
  ) => {
    if (!slug || !responseId) return;

    try {
      setIsLoading(true);
      const [, status] = await surveyResponseComplete(slug, responseId, values);

      if (isSuccess(status)) {
        setStep(3); // Success step
      } else {
        setError("Failed to submit survey. Please try again later.");
      }
    } catch (err) {
      setError("Failed to submit survey. Please try again later.");
      console.error("Error submitting survey:", err);
    } finally {
      setIsLoading(false);
      setSubmitting(false);
    }
  };

  const generateValidationSchema = () => {
    if (!survey) return {};

    const schemaShape: Record<string, any> = {};

    survey.questions.forEach((question) => {
      if (question.required) {
        switch (question.type) {
          case "text":
            schemaShape[question.id] = Yup.string().required(
              "This field is required"
            );
            break;
          case "number":
            schemaShape[question.id] = Yup.number().required(
              "This field is required"
            );
            break;
          case "select":
            schemaShape[question.id] = Yup.string().required(
              "This field is required"
            );
            break;
          case "multi-select":
            schemaShape[question.id] = Yup.array()
              .min(1, "At least one option must be selected")
              .required("This field is required");
            break;
          case "boolean":
            schemaShape[question.id] = Yup.boolean().required(
              "This field is required"
            );
            break;
          default:
            schemaShape[question.id] = Yup.mixed().required(
              "This field is required"
            );
        }
      }
    });

    return Yup.object().shape(schemaShape);
  };

  const generateInitialValues = () => {
    if (!survey) return {};

    const initialValues: Record<string, any> = {};

    survey.questions.forEach((question) => {
      switch (question.type) {
        case "text":
          initialValues[question.id] = "";
          break;
        case "number":
          initialValues[question.id] = "";
          break;
        case "select":
          initialValues[question.id] = "";
          break;
        case "multi-select":
          initialValues[question.id] = [];
          break;
        case "boolean":
          initialValues[question.id] = false;
          break;
        default:
          initialValues[question.id] = "";
      }
    });

    return initialValues;
  };

  const renderField = (question: any) => {
    switch (question.type) {
      case "text":
        return (
          <Field
            type="text"
            id={question.id}
            name={question.id}
            placeholder={question.placeholder || ""}
            className="tw:w-full tw:p-3 tw:border tw:border-gray-200 tw:rounded-lg"
          />
        );
      case "number":
        return (
          <Field
            type="number"
            id={question.id}
            name={question.id}
            placeholder={question.placeholder || ""}
            className="tw:w-full tw:p-3 tw:border tw:border-gray-200 tw:rounded-lg"
          />
        );
      case "select":
        return (
          <Field
            as="select"
            id={question.id}
            name={question.id}
            className="tw:w-full tw:p-3 tw:border tw:border-gray-200 tw:rounded-lg"
          >
            <option value="">Select...</option>
            {question.options?.map((option: any) => (
              <option key={option.id} value={option.id}>
                {option.value}
              </option>
            ))}
          </Field>
        );
      case "multi-select":
        return (
          <div className="checkbox-group">
            {question.options?.map((option: any) => (
              <div key={option.id} className="checkbox-item">
                <Field
                  type="checkbox"
                  id={`${question.id}-${option.id}`}
                  name={question.id}
                  value={option.id}
                  className="tw:w-5 tw:h-5 tw:mr-3 tw:text-primary tw:rounded"
                />
                <label htmlFor={`${question.id}-${option.id}`}>
                  {option.value}
                </label>
              </div>
            ))}
          </div>
        );
      case "boolean":
        return (
          <div className="checkbox-item">
            <Field
              type="checkbox"
              id={question.id}
              name={question.id}
              className="tw:w-5 tw:h-5 tw:mr-3 tw:text-primary tw:rounded"
            />
            <label htmlFor={question.id}>Yes</label>
          </div>
        );
      default:
        return <div>Unsupported question type</div>;
    }
  };

  if (isLoading && step === 1) {
    return (
      <div className="survey-participate-container">
        <div className="survey-participate-header">
          <h1>Loading Survey...</h1>
          <Link to="/ongoing-surveys" className="back-link">
            Back to Surveys
          </Link>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="survey-participate-container">
        <div className="survey-participate-header">
          <h1>Error</h1>
          <p>{error}</p>
          <Link to="/ongoing-surveys" className="back-link">
            Back to Surveys
          </Link>
        </div>
      </div>
    );
  }

  if (!survey) {
    return (
      <div className="survey-participate-container">
        <div className="survey-participate-header">
          <h1>Survey Not Found</h1>
          <Link to="/ongoing-surveys" className="back-link">
            Back to Surveys
          </Link>
        </div>
      </div>
    );
  }

  if (step === 3) {
    return (
      <div className="survey-participate-container">
        <div className="survey-participate-header">
          <h1>Thank You!</h1>
          <p>Your response has been recorded.</p>
          <Link to="/ongoing-surveys" className="back-link">
            Back to Surveys
          </Link>
        </div>
        <div className="success-message">
          <h2>Survey Completed Successfully</h2>
          <p>
            Thank you for participating in the survey. Your response will help
            contribute to salary transparency.
          </p>
          <Link to="/results" className="back-link">
            View Survey Results
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="survey-participate-container">
      <div className="survey-participate-header">
        <h1>{survey.title}</h1>
        <p>Participate in this salary survey to contribute to transparency</p>
        <Link to="/ongoing-surveys" className="back-link">
          Back to Surveys
        </Link>
      </div>

      {step === 1 && (
        <div className="survey-info-card">
          <h2>Survey Information</h2>
          <p className="survey-info-description">{survey.description}</p>

          <div className="survey-info-stats">
            <div className="survey-stat-item">
              <span className="survey-stat-label">Created At</span>
              <span className="survey-stat-value">
                {new Date(survey.createdAt).toLocaleDateString()}
              </span>
            </div>
            <div className="survey-stat-item">
              <span className="survey-stat-label">Last Updated</span>
              <span className="survey-stat-value">
                {new Date(survey.updatedAt).toLocaleDateString()}
              </span>
            </div>
            <div className="survey-stat-item">
              <span className="survey-stat-label">Min Completion Time</span>
              <span className="survey-stat-value">
                {survey.minCompletionTimeMin} minutes
              </span>
            </div>
          </div>

          <div className="survey-actions">
            <div className="anonymity-toggle">
              <label className="toggle-switch">
                <input
                  type="checkbox"
                  checked={isAnonymous}
                  onChange={() => setIsAnonymous(!isAnonymous)}
                />
                <span className="toggle-slider"></span>
              </label>
              <span>
                {isAnonymous
                  ? "Anonymous Participation"
                  : "Non-Anonymous Participation"}
              </span>
            </div>

            <button
              className="start-survey-button"
              onClick={startSurvey}
              disabled={isLoading}
            >
              {isLoading ? "Starting..." : "Start Survey"}
            </button>
          </div>
        </div>
      )}

      {step === 2 && (
        <div className="survey-form-container">
          <h2>Complete the Survey</h2>
          <Formik
            initialValues={generateInitialValues()}
            validationSchema={generateValidationSchema()}
            onSubmit={handleSubmit}
          >
            {({ isSubmitting }) => (
              <Form>
                {survey.questions
                  .sort((a, b) => a.order - b.order)
                  .map((question) => (
                    <div key={question.id} className="form-group">
                      <label htmlFor={question.id}>
                        {question.text}
                        {question.required && (
                          <span className="text-red-500">*</span>
                        )}
                      </label>
                      {renderField(question)}
                      <ErrorMessage
                        name={question.id}
                        component="div"
                        className="error-message"
                      />
                    </div>
                  ))}

                <div className="form-actions">
                  <button
                    type="button"
                    className="secondary-button tw:bg-gray-100 tw:text-gray-600 tw:px-6 tw:py-3 tw:rounded-lg tw:font-medium"
                    onClick={() => setStep(1)}
                  >
                    Back
                  </button>
                  <button
                    type="submit"
                    className="submit-button"
                    disabled={isSubmitting || isLoading}
                  >
                    {isSubmitting || isLoading
                      ? "Submitting..."
                      : "Submit Response"}
                  </button>
                </div>
              </Form>
            )}
          </Formik>
        </div>
      )}
    </div>
  );
};

export default SurveyParticipate;
