import { useState, useEffect } from 'react';
import { Link, useParams } from "react-router-dom";
import "../styles/Results.css";
import { isSurveyDetail, Survey } from "../types/survey.types";
import { surveyGet, surveyGetResults } from "../services/survey.api";
import { isSuccess } from "../services/base.api";
import {
  isResponseGraphListItemArray,
  ResponseGraphListItem,
} from "../types/response-graph.types";
import SurveyResultGraph from "../components/app/survey/SurveyResultGraph";

const Results = () => {
  const { slug } = useParams<{ slug: string }>();
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isResultsLoading, setIsResultsLoading] = useState(true);
  const [results, setResults] = useState<ResponseGraphListItem[]>([]);

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
          surveyGetResults(slug)
            .then(([results, status]) => {
              if (isSuccess(status) && isResponseGraphListItemArray(results)) {
                setResults(results);
              }
            })
            .finally(() => {
              setIsResultsLoading(false);
            });
          return;
        }
        setError("Failed to load survey. Please try again later.");
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [slug]);

  if (isLoading || isResultsLoading) {
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

  return (
    <div className="results-container loaded">
      <div className="results-header tw:shadow-lg tw:rounded-b-xl">
        <h1 className="tw:relative">
          Salary <span className="highlight tw:inline-block">Insights</span>
        </h1>
        <p className="subtitle tw:max-w-2xl">
          Explore comprehensive data on software engineering compensation trends
        </p>
        <Link to="/" className="back-link tw:transition-all tw:duration-300">
          <span className="back-arrow">←</span> Back to Home
        </Link>
      </div>

      <div className="dashboard-container">
        <div className="dashboard-section fade-in">
          <div className="section-header tw:py-6">
            <h2 className="tw:text-3xl tw:font-bold tw:mb-3">
              How Experience Impacts Salary
            </h2>
            <p className="tw:text-gray-600 tw:max-w-3xl tw:mx-auto">
              Analysis of the relationship between years of experience and
              compensation
            </p>
          </div>

          <div className="chart-grid">
            {results.map((result) => (
              <SurveyResultGraph key={result.id} {...result} />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Results; 