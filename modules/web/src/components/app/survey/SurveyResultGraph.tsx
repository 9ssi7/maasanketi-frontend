import {
  ResponseGraphKind,
  ResponseGraphListItem,
} from "../../../types/response-graph.types";
import BarGraphContent from "./graphs/BarGraphContent";
import LineGraphContent from "./graphs/LineGraphContent";

type Props = ResponseGraphListItem;

function GraphContent(props: Props) {
  switch (props.kind) {
    case ResponseGraphKind.Bar:
      return <BarGraphContent {...props} />;
    case ResponseGraphKind.Line:
      return <LineGraphContent {...props} />;
    default:
      return <></>;
  }
}

export default function SurveyResultGraph(props: Props) {
  return (
    <div className="chart-card">
      <h3>{props.title}</h3>
      <div className="chart-description">
        <p>{props.description}</p>
      </div>
      <GraphContent {...props} />
    </div>
  );
}
