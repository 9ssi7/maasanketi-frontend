import { ListResponse } from "../types/base.types";

export function listMerger<T = any>(
  prev: ListResponse<T> | null,
  next: ListResponse<T>
): ListResponse<T> {
  return {
    ...next,
    list: [...(prev?.list || []), ...next.list],
  };
}
