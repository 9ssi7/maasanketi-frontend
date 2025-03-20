import { Bar } from "recharts";

import { Legend } from "recharts";

import {
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { ResponseGraphListItem } from "../../../../types/response-graph.types";

export default function BarGraphContent(props: ResponseGraphListItem) {
  const data = props.content
    .filter((item) => !!item.values && !!item.labels)
    .flatMap((item) =>
      item.labels.map((label, index) => ({
        label,
        value: item.values[index] ? item.values[index] : "0",
      }))
    );
  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart
        data={data}
        margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
      >
        <CartesianGrid strokeDasharray="3 3" stroke="#eee" />
        <XAxis dataKey="label" />
        <YAxis />
        <Tooltip
          formatter={(value) => [
            `${value.toLocaleString()}`,
            props.graphContent.valueFields[0].label,
          ]}
          contentStyle={{
            borderRadius: "8px",
            border: "none",
            boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
          }}
        />
        <Legend />
        <Bar
          dataKey="value"
          fill="#4361ee"
          name={props.graphContent.valueFields[0].label}
          radius={[4, 4, 0, 0]}
        />
      </BarChart>
    </ResponsiveContainer>
  );
}
