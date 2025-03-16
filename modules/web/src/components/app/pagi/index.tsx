import AppIconChevronLeft from "../icon/chevron-left";
import AppIconChevronRight from "../icon/chevron-right";
import "./pagi.css";

type Props = {
  page: number;
  nextPossible: boolean;
  onChange: (page: number) => void;
};

const getPageNumbers = (page: number): number[] => {
  if (page < 1) return [1];
  if (page === 1) return [1];
  if (page <= 2) return [1, 2];
  if (page <= 3) return [1, 2, 3];
  return [page - 2, page - 1, page];
};

export default function Pagi({ page, nextPossible, onChange }: Props) {
  const numbers = getPageNumbers(page);
  return (
    <div className="pagination">
      <a
        className={`pagination-button ${
          page <= 1 || numbers.some((n) => n === 1) ? "disabled" : ""
        }`}
        onClick={() => onChange(page - 1)}
      >
        <AppIconChevronLeft />
      </a>
      {numbers.map((n) => (
        <a
          key={n}
          className={`pagination-button ${n === page ? "active" : ""}`}
          onClick={() => onChange(n)}
        >
          {n}
        </a>
      ))}
      <a
        className={`pagination-button ${!nextPossible ? "disabled" : ""}`}
        onClick={() => onChange(page + 1)}
      >
        <AppIconChevronRight />
      </a>
    </div>
  );
}
