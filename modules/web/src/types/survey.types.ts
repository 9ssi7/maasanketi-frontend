// Question types
export type QuestionType = 'text' | 'number' | 'select' | 'multi-select' | 'boolean';

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
  answers: Record<string, any>;
  startedAt: string;
  completedAt?: string;
  ipAddress: string;
  userAgent: string;
  isCompleted: boolean;
  isAnonymous: boolean;
  createdAt: string;
  updatedAt: string;
}

// API request/response types

// Response for listing surveys
export interface ListSurveysResponse {
  surveys: Survey[];
}

// Response for getting a survey
export type GetSurveyResponse = Survey;

// Request for starting a survey response - no body needed
export type StartSurveyResponseRequest = Record<string, never>;

// Response for starting a survey response
export type StartSurveyResponseResponse = SurveyResponse;

// Request for updating a survey response
export interface UpdateSurveyResponseRequest {
  answers: Record<string, any>;
}

// Response for updating a survey response
export type UpdateSurveyResponseResponse = SurveyResponse;

// Request for completing a survey response - no body needed
export type CompleteSurveyResponseRequest = Record<string, never>;

// Response for completing a survey response
export type CompleteSurveyResponseResponse = SurveyResponse;

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

export function isSurveyListItem(survey: unknown): survey is SurveyListItem {
  return typeof survey === 'object' && survey !== null && 'id' in survey
}

export function isSurveyList(surveys: unknown): surveys is SurveyListItem[] {
  return Array.isArray(surveys) && surveys.every(isSurveyListItem)
}