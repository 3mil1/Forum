import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'

// initialize an empty api service that we'll inject endpoints into later as needed
export const api = createApi({
    reducerPath: 'api',
    baseQuery: fetchBaseQuery({baseUrl: 'http://localhost:8081/api', credentials: "include"}),
    tagTypes: ['Users', 'Auth', 'Mark', 'Post'],
    endpoints: () => ({}),
})