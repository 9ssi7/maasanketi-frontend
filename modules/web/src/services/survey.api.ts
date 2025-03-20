import { ListResponse } from './../types/base.types';
import { SurveyAnswers, SurveyDetail, SurveyListItem, SurveyResponseStartResponse } from "../types/survey.types"
import { req, withQuery } from "./base.api"
import { ResponseGraphListItem } from '../types/response-graph.types';

export const surveyList = async(q?: string) => {
    return req<ListResponse<SurveyListItem>>(withQuery("api/v1/surveys", q))
}

export const surveyGet = async(slug: string) => {
    return req<SurveyDetail>(`api/v1/surveys/${slug}`)
}

export const surveyResponseStart = async(slug: string, isAnonymous: boolean) => {
    return req<SurveyResponseStartResponse>(`api/v1/surveys/${slug}/responses`, {
        method: "POST",
        body: JSON.stringify({ isAnonymous })
    })
}

export const surveyResponseComplete = async(slug: string, responseId: string, answers: SurveyAnswers) => {
    return req(`api/v1/surveys/${slug}/responses/${responseId}`, {
        method: "PATCH",
        body: JSON.stringify({ answers })
    })
}

export const surveyGetResults = async(slug: string) => {
    return req<ResponseGraphListItem[]>(`api/v1/surveys/${slug}/results`)
}