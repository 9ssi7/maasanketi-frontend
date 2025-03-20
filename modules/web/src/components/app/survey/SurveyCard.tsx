import { Link } from "react-router-dom";
import { SurveyListItem } from "../../../types/survey.types";

type Props = SurveyListItem;

export default function SurveyCard({ ...survey }: Props) {
  return (
    <div key={survey.id} className="survey-card">
      <h2>{survey.title}</h2>
      <p className="survey-description">{survey.description}</p>

      <div className="survey-stats">
        <div className="stat-item">
          <span className="stat-label">Participants:</span>
          <span className="stat-value">{survey.participants}</span>
        </div>
        {!survey.isExpired && (
          <div className="stat-item">
            <span className="stat-label">Closing:</span>
            <span className="stat-value">
              {new Date(survey.updatedAt).toLocaleDateString()}
            </span>
          </div>
        )}
        {survey.isExpired && (
          <>
            <div className="stat-item">
              <span className="stat-label">Expired:</span>
              <span className="stat-value">Yes</span>
            </div>
            <div className="stat-item">
              <span className="stat-label">Closed at:</span>
              <span className="stat-value">
                {new Date(survey.updatedAt).toLocaleDateString()}
              </span>
            </div>
          </>
        )}
        {survey.createdBy && (
          <div className="stat-item">
            <span className="stat-label">Created by:</span>
            <span className="stat-value">{survey.createdBy}</span>
          </div>
        )}
      </div>

      <div className="survey-actions">
        {!survey.isExpired && (
          <Link
            to={`/surveys/${survey.slug}/participate`}
            className="participate-button"
          >
            Participate
          </Link>
        )}
        {survey.isExpired && (
          <Link
            to={`/surveys/${survey.slug}/results`}
            className="participate-button"
          >
            View Results
          </Link>
        )}
      </div>
    </div>
  );
}
