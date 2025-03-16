export interface ListResponse<Entity> {
	list: Entity[];
	limit: number;
	page: number;
}

export function isListResponse(arg: any): arg is ListResponse<any> {
	return (
		!!arg &&
		arg.list !== undefined &&
		arg.page !== undefined &&
		Array.isArray(arg.list)
	);
}

export type PropsWithClassName<T = unknown> = T &{
	className?: string;
};
