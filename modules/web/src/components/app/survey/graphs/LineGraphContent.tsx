import {
  CartesianGrid,
  Legend,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { ResponseGraphListItem } from "../../../../types/response-graph.types";

export default function LineGraphContent(props: ResponseGraphListItem) {
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
      <LineChart
        data={data}
        margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
      >
        <CartesianGrid strokeDasharray="3 3" stroke="#eee" />
        <XAxis dataKey="label" />
        <YAxis />
        <Tooltip
          formatter={(value) => [
            `$${value.toLocaleString()}`,
            "Average Salary",
          ]}
          contentStyle={{
            borderRadius: "8px",
            border: "none",
            boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
          }}
        />
        <Legend />
        <Line
          type="monotone"
          dataKey="value"
          stroke="#4361ee"
          strokeWidth={3}
          activeDot={{ r: 8 }}
          name="Average Salary ($)"
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
