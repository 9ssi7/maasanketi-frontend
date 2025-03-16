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
        <div className="stat-item">
          <span className="stat-label">Closing:</span>
          <span className="stat-value">
            {new Date(survey.updatedAt).toLocaleDateString()}
          </span>
        </div>
        {survey.createdBy && (
          <div className="stat-item">
            <span className="stat-label">Created by:</span>
            <span className="stat-value">{survey.createdBy}</span>
          </div>
        )}
      </div>

      <div className="survey-actions">
        <button className="participate-button">Participate</button>
        <button className="details-button">View Details</button>
      </div>
    </div>
  );
}
