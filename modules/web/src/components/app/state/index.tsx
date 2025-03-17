import classNames from "classnames";
import Loader from "../loader";
import AppIconError from "../icon/error";
import AppIconXCircle from "../icon/x-circle";
import AppIconSearch from "../icon/search";

type ErrorProps = {
  data?: any;
};

type EmptyProps = {
  filtered?: boolean;
  text?: string;
  onClear?: () => void;
  class?: string;
};

function State() {
  return <></>;
}

type LoadingProps = {
  className?: string;
};

State.Loading = ({ className }: LoadingProps) => (
  <div
    className={classNames(
      "tw:p-4 tw:flex tw:items-center tw:justify-center",
      className
    )}
  >
    <Loader />
  </div>
);

State.Error = ({ data }: ErrorProps) => {
  return (
    <div className="tw:flex tw:flex-col tw:gap-1 tw:items-center tw:justify-center">
      <AppIconError className="tw:size-8 tw:text-red-500" />
      <div className="tw:text-lg tw:text-red-600">Hata</div>
      <p className="tw:text-sm tw:text-red-700">
        {data && typeof data === "object" && "message" in data
          ? data.message
          : "Bir hata oluştu"}
      </p>
    </div>
  );
};

State.Empty = ({
  filtered,
  onClear,
  text = "Sonuç Yok",
  class: className,
}: EmptyProps) => {
  return (
    <div
      className={classNames(
        className,
        "tw:flex tw:flex-col tw:gap-1 tw:items-center tw:justify-center tw:bg-white tw:rounded-lg tw:!py-10 tw:px-2"
      )}
    >
      <AppIconSearch className="tw:size-10 tw:text-neutral-700" />
      <div className="tw:text-lg tw:font-semibold tw:text-neutral-600">
        {text}
      </div>
      {filtered && (
        <p className="tw:text-sm tw:text-neutral-700">
          Filtreleme sonucunda hiçbir sonuç bulunamadı.
        </p>
      )}
      {filtered && onClear && (
        <button
          onClick={onClear}
          className="tw:text-sm tw:!mt-2 tw:text-white tw:hover:underline"
        >
          Filtreyi Temizle
        </button>
      )}
    </div>
  );
};

State.NotFound = () => (
  <div className="tw:flex tw:flex-col tw:gap-1 tw:items-center tw:justify-center">
    <AppIconSearch className="tw:size-8 tw:text-neutral-500" />
    <div className="tw:text-lg tw:text-neutral-600">Bulunamadı</div>
    <p className="tw:text-sm tw:text-neutral-700">
      Aradığınız içerik bulunamadı. Lütfen tekrar deneyin.
    </p>
  </div>
);

export default State;
