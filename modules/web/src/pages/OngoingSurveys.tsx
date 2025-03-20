import { Link, useSearchParams } from "react-router-dom";
import "../styles/OngoingSurveys.css";
import OngoingSurveySection from "../partials/OngoingSurveySection";
import { useRef, useState } from "react";
import {
  SurveyCategoryTexts,
  SurveySort,
  SurveySortTexts,
} from "../types/survey.types";

const OngoingSurveys = () => {
  const [searchParams] = useSearchParams();
  const [sort, setSort] = useState<SurveySort>(SurveySort.CreatedAtDesc);
  const [tag, setTag] = useState<string>("");
  const [search, setSearch] = useState<string>("");
  const [hideExpired, setHideExpired] = useState<boolean>(
    searchParams.has("hideExpired")
      ? searchParams.get("hideExpired") === "true"
      : true
  );
  const searchRef = useRef<HTMLInputElement>(null);
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
          <input
            id="search"
            type="text"
            placeholder="Anket Ara..."
            ref={searchRef}
          />
          <button onClick={() => setSearch(searchRef.current?.value || "")}>
            Ara
          </button>
        </div>

        <div className="filter-options">
          <select defaultValue={tag} onChange={(e) => setTag(e.target.value)}>
            <option value="">Tüm Departmanlar</option>
            {Object.keys(SurveyCategoryTexts).map((category) => (
              <option key={category} value={category}>
                {SurveyCategoryTexts[category]}
              </option>
            ))}
          </select>

          <select
            defaultValue={sort}
            onChange={(e) => setSort(e.target.value as SurveySort)}
          >
            {Object.values(SurveySort).map((sort) => (
              <option key={sort} value={sort}>
                {SurveySortTexts[sort]}
              </option>
            ))}
          </select>
        </div>
      </div>

      <OngoingSurveySection
        sort={sort}
        tag={tag}
        search={search}
        hideExpired={hideExpired}
        onClear={() => {
          setSearch("");
          setTag("");
          setSort(SurveySort.CreatedAtDesc);
          setHideExpired(true);
        }}
      />
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
