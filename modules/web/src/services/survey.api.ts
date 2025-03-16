import { ListResponse } from './../types/base.types';
import { Survey, SurveyListItem } from "../types/survey.types"
import { req, withQuery } from "./base.api"

export const surveyList = async(q?: string) => {
    return req<ListResponse<SurveyListItem>>(withQuery("api/v1/surveys", q))
}

export const surveyGet = async(slug: string) => {
    return req<Survey>(`api/v1/surveys/${slug}`)
}

