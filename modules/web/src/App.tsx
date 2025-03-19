import { Routes, Route } from 'react-router-dom';
import Landing from './pages/Landing';
import Results from './pages/Results';
import CreateSurvey from './pages/CreateSurvey';
import OngoingSurveys from './pages/OngoingSurveys';
import SurveyParticipate from "./pages/SurveyParticipate";
import "./App.css";

function App() {
  return (
    <div className="app">
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/results" element={<Results />} />
        <Route path="/create-survey" element={<CreateSurvey />} />
        <Route path="/ongoing-surveys" element={<OngoingSurveys />} />
        <Route
          path="/surveys/:slug/participate"
          element={<SurveyParticipate />}
        />
      </Routes>
    </div>
  );
}

export default App;
