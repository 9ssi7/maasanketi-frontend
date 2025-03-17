import classNames from "classnames";
import "./index.css";

type Size = "xs" | "sm" | "md" | "lg" | "xl";

type Props = {
  size?: Size;
};

const sizes = {
  xs: "tw:size-4",
  sm: "tw:size-5",
  md: "tw:size-6",
  lg: "tw:size-7",
  xl: "tw:size-8",
};

export default function Loader({ size = "md" }: Props) {
  return <div className={classNames("loader", sizes[size])} />;
}
