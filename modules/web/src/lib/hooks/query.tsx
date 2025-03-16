import { useEffect, useState } from "react";
import { ReqResult } from "../../services/base.api";
import { isApiError } from "../../services/error";
import { isListResponse } from "../../types/base.types";

type Fetcher<T = any, P = any> = (p?: P) => ReqResult<T>;

type Params<T, P> = {
  fetcher: Fetcher<T, P>;
  params?: P;
  dataMerger?: (prev: T | null, next: T) => T;
};

type Result<T, P> = {
  data: T | null;
  isLoading: boolean;
  isError: boolean;
  isNext: boolean;
  error: any;
  refetch: (p?: P) => void;
  setData: (data: T) => void;
};

export function useQuery<T = any, P = any>({
  fetcher,
  params,
  dataMerger,
}: Params<T, P>): Result<T, P> {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);
  const [isNext, setIsNext] = useState(false);
  const [error, setError] = useState<any>(null);

  useEffect(() => {
    const fetchData = async () => {
      setIsError(false);
      setIsLoading(true);

      const res = await fetcher(params);
      if (!res) {
        setIsError(true);
        setIsLoading(false);
        return;
      }

      const [data, status] = res;
      if (status !== 200 || isApiError(data)) {
        setIsError(true);
        setError(data);
        setIsLoading(false);
        return;
      }
      setData(data);
      if (isListResponse(data)) {
        setIsNext(data.limit === data.list?.length);
      }
      setIsLoading(false);
    };

    fetchData();
  }, []);

  const refetch = async (p?: P) => {
    setIsError(false);
    setIsLoading(true);

    const res = await fetcher(p || params);
    if (!res) {
      setIsError(true);
      setIsLoading(false);
      return;
    }

    const [data, status] = res;
    if (status !== 200 || isApiError(data)) {
      setIsError(true);
      setError(data);
      setIsLoading(false);
      return;
    }
    if (dataMerger) {
      setData((prev) => dataMerger(prev, data));
    } else {
      setData(data);
    }
    if (isListResponse(data)) {
      setIsNext((data as any).limit === (data as any).list.length);
    }
    setIsLoading(false);
  };

  return { data, isLoading, isNext, isError, error, refetch, setData };
}
