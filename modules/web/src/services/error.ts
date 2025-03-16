import { toast } from "sonner";
import { isSuccess } from "./base.api";

export type AuthRequiredError = object;

export type NotFoundError = object;

export type TooManyRequestsError = {
	reset?: number;
	remaining?: number;
	limit?: number;
	retryAfter?: number;
};

export type ValidationError = {
	data: {
		field: string;
		message: string;
		namespace: string;
		value: string;
	}[];
};

export type ApiError =
	| AuthRequiredError
	| TooManyRequestsError
	| ValidationError
	| NotFoundError;

export function isValidationError(error: any): error is ValidationError {
	return (
		error && Array.isArray(error.data) && error.data[0].field !== undefined
	);
}

export function isTooManyRequestsError(
	error: any,
): error is TooManyRequestsError {
	return error?.retryAfter !== undefined;
}
export function isApiError(error: any): error is ApiError {
	return isValidationError(error) || isTooManyRequestsError(error);
}

type Form = {
	setFieldError: (field: string, value: string | undefined) => void;
};

type Opts = {
	form?: Form;
	noToast?: boolean;
	onErrorField?: (f: string) => void;
};

export const handleApiErrorResult = (res: any, opts?: Opts) => {
	if (isValidationError(res)) {
		for (const error of res.data) {
			opts?.form?.setFieldError(error.field?.toLowerCase(), error.message);
			opts?.onErrorField?.(error.field?.toLowerCase());
		}
	}
	if (!!res && res.message) {
		if (!opts?.noToast) {
			toast.error(res.message);
		}
		return;
	}
	if (!!res && res.error) {
		if (!opts?.noToast) {
			toast.error(res.error);
		}
		return;
	}
};

const errorStates = [
	{
		code: 429,
		text: "Çok fazla istek gönderdiniz. Lütfen biraz bekleyin",
	},
	{
		code: 401,
		text: "Bu işlemi yapabilmek için giriş yapmalısınız",
	},
	{
		code: 403,
		text: "Bu işlemi yapabilmek için yetkiniz yok",
	},
	{
		code: 500,
		text: "Sunucu hatası oluştu",
	},
	{
		code: 422,
		text: "Form doğrulaması başarısız oldu",
	},
] as const;

export const validateWithStatus = (data: any, status: number) => {
	if (isSuccess(status)) return false;
	if (data === null) {
		toast.error("Bir hata oluştu");
		return true;
	}
	if (typeof data === "object" && data.message) {
		toast.error(data.message);
		return true;
	}
	const state = errorStates.find((e) => e.code === status);
	toast.error(state ? state.text : "Bir hata oluştu");
	return true;
};