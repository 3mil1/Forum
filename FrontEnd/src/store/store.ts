import {combineReducers, configureStore} from '@reduxjs/toolkit';
import {Reducer} from 'redux';
import {RESET_STATE_ACTION_TYPE} from './actions';
import {unauthenticatedMiddleware} from './middleware';
import {api} from "../services/api";

const reducers = {
    [api.reducerPath]: api.reducer,
};

const combinedReducer = combineReducers<typeof reducers>(reducers);

export const rootReducer: Reducer<RootState> = (state, action) => {
    if (action.type === RESET_STATE_ACTION_TYPE) {
        state = {} as RootState;
    }

    return combinedReducer(state, action);
}

export const setupStore = () => {
    return configureStore({
        reducer: rootReducer,
        middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat([
            // unauthenticatedMiddleware,
            api.middleware
        ]),
    });
}

export type RootState = ReturnType<typeof combinedReducer>;
export type AppStore = ReturnType<typeof setupStore>
export type AppDispatch = AppStore['dispatch']
