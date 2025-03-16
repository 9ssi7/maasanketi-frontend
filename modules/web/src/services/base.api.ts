import type { ApiError, NotFoundError, TooManyRequestsError } from "./error";

export const ApiUrl = import.meta.env.VITE_BASE_API_URL;

type StateHandler = (response: Response) => [ApiError, number];

const tooManyRequestsStateHandler: StateHandler = (
	response: Response,
): [TooManyRequestsError, 429] => {
	const resp: TooManyRequestsError = {
		limit: Number.parseInt(response.headers.get("X-Ratelimit-Limit") || "0"),
		remaining: Number.parseInt(
			response.headers.get("X-Ratelimit-Remaining") || "0",
		),
		reset: Number.parseInt(response.headers.get("X-Ratelimit-Reset") || "0"),
		retryAfter: Number.parseInt(response.headers.get("Retry-After") || "0"),
	};
	return [resp, 429];
};

const notFoundStateHandler: StateHandler = (): [ApiError, 404] => {
	return [{} as NotFoundError, 404];
};

const StateHandlers: Record<number, StateHandler> = {
	429: tooManyRequestsStateHandler,
	404: notFoundStateHandler,
};

export type ReqResult<T> = Promise<[T | ApiError, number]>;

export async function req<T = unknown>(
	endpoint: string,
	options: RequestInit = {},
): ReqResult<T> {
	const extraHeads: Record<string, string> = {};
	if (!(options.body instanceof FormData)) {
		extraHeads["Content-Type"] = "application/json";
	}
	const response = await fetch(`${ApiUrl}/${endpoint}`, {
		...options,
		//credentials: "include",

		headers: {
			...options.headers,
			...extraHeads,
		},
	}).catch(() => {
		return new Response(null, { status: 500 });
	});
	if (!response.ok) {
		const handler = StateHandlers[response.status];
		if (handler) {
			return handler(response);
		}
		try {
			const json = await response.json();
			return [json, response.status];
		} catch (e: unknown) {
			console.error(e);
        }
		return [{} as ApiError, response.status];
	}
	if (response.headers.get("Content-Type")?.includes("application/json")) {
		return [await response.json(), response.status];
	}
	return [await response.text() as T, response.status];
}

export const withQuery = (endpoint: string, query?: string): string => {
	return query ? `${endpoint}?${query}` : endpoint;
};

export const jsonToQuery = (json: Record<string, unknown>): string => {
	if (!json) return "";
	return Object.entries(json)
		.filter(([, value]) => value !== undefined && value !== null)
		.map(([key, value]) => `${key}=${value}`)
		.join("&");
};

export const queryToJson = (query: string): Record<string, string> => {
	if (!query) return {};
	const params = new URLSearchParams(query);
	const json: Record<string, string> = {};
	for (const [key, value] of params.entries()) {
		json[key] = value;
	}
	return json;
};

type FormDataOpts = {
	arrayWithBrackets?: boolean;
	arrayWithBracketAndIndex?: boolean;
};

export const jsonToFormData = (
	data: { [s: string]: unknown },
	{
		arrayWithBrackets = false,
		arrayWithBracketAndIndex = true,
	}: FormDataOpts = {},
): FormData => {
	const formData = new FormData();
	for (const [key, value] of Object.entries(data)) {
		if (Array.isArray(value)) {
			for (const idx in value) {
				if (arrayWithBrackets) {
					formData.append(`${key}[]`, value[idx]);
				} else if (arrayWithBracketAndIndex) {
					formData.append(`${key}[${idx}]`, value[idx]);
				} else {
					formData.append(key, value[idx]);
				}
			}
			continue;
		}
		if (typeof value !== "undefined" && value !== null) {
			if (typeof value === "string" || value instanceof Blob) {
				formData.append(key, value);
			} else {
				formData.append(key, JSON.stringify(value));
			}
		}
	}
	return formData;
};

export const isSuccess = (status: number): boolean => {
	return status >= 200 && status < 300;
};