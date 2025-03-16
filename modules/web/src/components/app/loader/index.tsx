import classNames from "classnames";
import "./index.css";

type Size = "xs" | "sm" | "md" | "lg" | "xl";

type Props = {
  size?: Size;
};

const sizes = {
  xs: "size-4",
  sm: "size-5",
  md: "size-6",
  lg: "size-7",
  xl: "size-8",
};

export default function Loader({ size = "md" }: Props) {
  return <div className={classNames("loader", sizes[size])} />;
}
