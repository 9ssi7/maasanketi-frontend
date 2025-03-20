export enum ResponseGraphKind {
    Bar = "bar",
    Line = "line",
    Pie = "pie",
    Radar = "radar",
    Scatter = "scatter",
}

export type ResponseGraphListItem = {
    id: string;
    surveyId: string;
    graphId: string;
    title: string;
    description: string;
    kind: ResponseGraphKind;
    content: ResponseGraphContent[];
    graphContent: GraphContent;
    createdAt: string;
    updatedAt: string;
}

export type ResponseGraphContent = {
    labels: string[];
    values: string[];
}

export type GraphContent = {
    keyFields: KeyField[];
    valueFields: ValueField[];
}

export type KeyField = {
    field: string;
    label: string;
}

export type ValueField = {
    field: string;
    label: string;
    strategy: ValueStrategy;
}

export enum ValueStrategy {
    Count = "count",
    Sum = "sum",
    Average = "average",
}

export function isResponseGraphListItem(item: unknown): item is ResponseGraphListItem {
    return typeof item === "object" && item !== null && "id" in item
}

export function isResponseGraphListItemArray(items: unknown): items is ResponseGraphListItem[] {
    return Array.isArray(items) && items.every(isResponseGraphListItem)
}