// Question types
export type QuestionType = 'text' | 'number' | 'select' | 'multi-select' | 'boolean';
export type SurveyAnswers = Record<string, any>;

// Option for select and multi-select questions
export interface Option {
  id: string;
  value: string;
}

// Conditional for showing a question based on another question's answer
export interface Conditional {
  questionId: string;
  value: string;
}

// Survey question
export interface Question {
  id: string;
  text: string;
  type: QuestionType;
  required: boolean;
  options?: Option[];
  placeholder?: string;
  conditional?: Conditional;
  order: number;
}

// Survey definition
export interface Survey {
  id: string;
  slug: string;
  title: string;
  description: string;
  questions: Question[];
  createdAt: string;
  updatedAt: string;
  minCompletionTimeMin: number;
}

// Survey response
export interface SurveyResponse {
  id: string;
  surveySlug: string;
  userId?: string;
  answers: SurveyAnswers;
  startedAt: string;
  completedAt?: string;
  ipAddress: string;
  userAgent: string;
  isCompleted: boolean;
  isAnonymous: boolean;
  createdAt: string;
  updatedAt: string;
}


// Response for getting a survey
export type GetSurveyResponse = Survey;

// Response for listing survey responses
export interface ListSurveyResponsesResponse {
  responses: SurveyResponse[];
}

// Response for getting survey statistics
export interface SurveyStats {
  totalResponses: number;
  completedResponses: number;
  averageCompletionTimeMinutes: number;
  questionStats: Record<string, {
    totalAnswers: number;
    answerDistribution: Record<string, number>;
  }>;
}

export interface SurveyListItem {
  id: string;
  slug: string;
  title: string;
  description: string;
  createdAt: string;
  updatedAt: string;
  minCompletionTimeMin: number;
  createdBy?: string;
  tags?: string[];
  participants: number;
}

export enum SurveySort {
  CreatedAtAsc = "created_at_asc",
  CreatedAtDesc = "created_at_desc",
  MostParticipants = "most_participants",
  FinishesAtAsc = "finishes_at_asc",
  FinishesAtDesc = "finishes_at_desc",
}

export interface SurveyListRequest {
  tag?: string;
  sort?: SurveySort;
  search?: string;
  hideExpired?: boolean;
}

export interface SurveyResponseStartRequest {
  surveySlug: string
  is_anonymous: boolean
}

export interface SurveyResponseStartResponse {
  id: string
  surveySlug: string
  startedAt: string
  answers: SurveyAnswers
}

export interface SurveyResponseCompleteRequest {
  slug: string
  responseId: string
  answers: SurveyAnswers
}

export function isSurveyListItem(survey: unknown): survey is SurveyListItem {
  return typeof survey === 'object' && survey !== null && 'id' in survey
}

export function isSurveyList(surveys: unknown): surveys is SurveyListItem[] {
  return Array.isArray(surveys) && surveys.every(isSurveyListItem)
}

export const SurveySortTexts : Record<SurveySort, string> = {
  [SurveySort.MostParticipants]: "En Fazla Katılımcı",
  [SurveySort.CreatedAtAsc]: "En Eski Oluşturulan",
  [SurveySort.CreatedAtDesc]: "En Yeni Oluşturulan",
  [SurveySort.FinishesAtAsc]: "En Yakın Bitiş",
  [SurveySort.FinishesAtDesc]: "En Uzak Bitiş",
}

export const SurveyCategoryTexts : Record<string, string> = {
  "frontend": "Frontend",
  "backend": "Backend",
  "fullstack": "Fullstack",
  "devops": "DevOps",
  "qa": "QA",
  "design": "Design",
  "pm": "Product Management",
  "hr": "Human Resources",
  "marketing": "Marketing",
  "sales": "Sales",
  "other": "Diğer",
}
