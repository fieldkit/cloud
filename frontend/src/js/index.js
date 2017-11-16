import '../scss/app.scss'

import 'babel-polyfill'
import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'
import { rootSaga } from './actions/sagas'
import createSagaMiddleware from 'redux-saga'
import * as actions from './actions'
import reducer from './reducers'

import RootContainer from './containers/Root/Root'
import MapPageContainer from './containers/MapPage/MapPage'

const sagaMiddleware = createSagaMiddleware()

import * as ActionTypes from './actions/types'

const loggerMiddleware = createLogger({
    predicate: (getState, action) => action.type != ActionTypes.SET_MOUSE_POSITION && action.type != ActionTypes.UPDATE_DATE, // __DEV__,
    collapsed: (getState, action) => true
});

const createStoreWithMiddleware = applyMiddleware(
    thunkMiddleware,
    sagaMiddleware,
    loggerMiddleware,
)(createStore)

const store = createStoreWithMiddleware(reducer)

sagaMiddleware.run(rootSaga);

const routes = (
    <Route path="/" component={RootContainer}>
        <IndexRoute />
        <Route path=":expeditionSlug">
            <IndexRoute component={MapPageContainer} />
            <Route path="map" component={MapPageContainer} />
        </Route>
    </Route>
)

const render = function() {
    ReactDOM.render(
        (
            <Provider store={ store }>
                <Router history={ browserHistory } routes={ routes } />
            </Provider>
        ),
        document.getElementById('fieldkit')
    )
}

render()
