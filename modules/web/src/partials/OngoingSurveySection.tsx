import { useState } from "react";
import Pagi from "../components/app/pagi";
import State from "../components/app/state";
import SurveyCard from "../components/app/survey/SurveyCard";
import { useQuery } from "../lib/hooks/query";
import { listMerger } from "../lib/pagi";
import { jsonToQuery } from "../services/base.api";
import { surveyList } from "../services/survey.api";
import { SurveySort } from "../types/survey.types";

type Props = {
  tag?: string;
  sort?: SurveySort;
  search?: string;
  hideExpired?: boolean;
  onClear?: () => void;
};

export default function OngoingSurveySection({
  tag,
  sort,
  hideExpired,
  search,
  onClear,
}: Props) {
  const [page, setPage] = useState(1);
  const { data, isLoading, isError } = useQuery(
    {
      fetcher: () =>
        surveyList(jsonToQuery({ tag, sort, hideExpired, search, page })),
      dataMerger: listMerger,
    },
    [tag, sort, hideExpired, search, page]
  );

  if (isLoading) return <State.Loading />;

  if (isError) return <State.Error data={data} />;

  return (
    <>
      <div className="survey-list">
        {data && data?.list?.length > 0 ? (
          data?.list.map((survey) => <SurveyCard key={survey.id} {...survey} />)
        ) : (
          <State.Empty
            filtered={
              !!search ||
              !!tag ||
              sort !== SurveySort.CreatedAtDesc ||
              hideExpired !== true
            }
            onClear={onClear}
          />
        )}
      </div>
      <Pagi
        page={page}
        nextPossible={data?.list?.length === data?.limit || false}
        onChange={setPage}
      />
    </>
  );
}
