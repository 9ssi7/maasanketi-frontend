import { useState, useEffect } from 'react';
import { Link, useParams } from "react-router-dom";
import "../styles/Results.css";
import SurveyResultScreen from "../components/app/survey/SurveyResultScreen";
import { isSurveyDetail, Survey } from "../types/survey.types";
import { surveyGet } from "../services/survey.api";
import { isSuccess } from "../services/base.api";

const Results = () => {
  const { slug } = useParams<{ slug: string }>();
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!slug) return;
    surveyGet(slug)
      .then(([survey, status]) => {
        if (isSuccess(status) && isSurveyDetail(survey)) {
          if (!survey.isExpired) {
            setError("Survey is not expired. Please try again later.");
            return;
          }
          setSurvey(survey);
          return;
        }
        setError("Failed to load survey. Please try again later.");
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [slug]);

  if (isLoading) {
    return (
      <div className="survey-participate-container">
        <div className="survey-participate-header">
          <h1>Loading Survey...</h1>
          <Link to="/surveys" className="back-link">
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
          <Link to="/surveys" className="back-link">
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
          <Link to="/surveys" className="back-link">
            Back to Surveys
          </Link>
        </div>
      </div>
    );
  }

  return <SurveyResultScreen />;
};

export default Results; 