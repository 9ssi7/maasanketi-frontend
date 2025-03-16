import { Link } from "react-router-dom";
import "../styles/OngoingSurveys.css";
import OngoingSurveySection from "../partials/OngoingSurveySection";

const OngoingSurveys = () => {
  return (
    <div className="ongoing-surveys-container">
      <div className="ongoing-surveys-header">
        <h1>Ongoing Salary Surveys</h1>
        <p>
          Participate in active surveys and contribute to salary transparency
        </p>
        <Link to="/" className="back-link">
          Back to Home
        </Link>
      </div>

      <div className="survey-filters">
        <div className="search-bar">
          <input type="text" placeholder="Search surveys..." />
          <button>Search</button>
        </div>

        <div className="filter-options">
          <select defaultValue="">
            <option value="">All Roles</option>
            <option value="frontend">Frontend</option>
            <option value="backend">Backend</option>
            <option value="fullstack">Fullstack</option>
            <option value="devops">DevOps</option>
          </select>

          <select defaultValue="newest">
            <option value="newest">Newest First</option>
            <option value="closing">Closing Soon</option>
            <option value="popular">Most Participants</option>
          </select>
        </div>
      </div>

      <OngoingSurveySection />

      <div className="ongoing-surveys-note">
        <p>
          All survey responses are anonymous and will be used only for
          statistical purposes.
        </p>
      </div>
    </div>
  );
};

export default OngoingSurveys;
