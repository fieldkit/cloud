import '../scss/app.scss'

import 'babel-polyfill'
import React from 'react'
import ReactDOM from 'react-dom'
import { createStore, applyMiddleware, combineReducers } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import multiMiddleware from 'redux-multi'
import { createLogger } from 'redux-logger'
import { Router, Route, IndexRoute, Redirect, browserHistory } from 'react-router'
import { rootSaga } from './actions/sagas'
import createSagaMiddleware from 'redux-saga'
import * as actions from './actions'
import reducer from './reducers'

import RootContainer from './containers/Root/Root'
import MapPageContainer from './containers/MapPage/MapPage'
import JournalPageContainer from './containers/JournalPage/JournalPage'
import DataPageContainer from './containers/DataPage/DataPage'

const sagaMiddleware = createSagaMiddleware()

import * as ActionTypes from './actions/types'

const loggerMiddleware = createLogger({
    predicate: (getState, action) => action.type != ActionTypes.SET_MOUSE_POSITION && action.type != ActionTypes.UPDATE_DATE, // __DEV__,
    collapsed: (getState, action) => true
});

const createStoreWithMiddleware = applyMiddleware(
    thunkMiddleware,
    multiMiddleware,
    sagaMiddleware,
    loggerMiddleware,
)(createStore)

const store = createStoreWithMiddleware(reducer)

sagaMiddleware.run(rootSaga);

const routes = (
<Route path="/" component={ RootContainer }>
    <IndexRoute />
    <Route path=":expeditionID">
        <IndexRoute component={ MapPageContainer } onEnter={ (state, replace) => {
                                                                 store.dispatch(actions.setCurrentPage('map'))
                                                             } } />
        <Route path="map" component={ MapPageContainer } onEnter={ (state, replace) => {
                                                                       store.dispatch(actions.setCurrentPage('map'))
                                                                   } } />
        <Route path="journal" component={ JournalPageContainer } onEnter={ (state, replace) => {
                                                                               store.dispatch(actions.selectPlaybackMode('pause'))
                                                                               store.dispatch(actions.selectFocusType('expedition'))
                                                                               store.dispatch(actions.setCurrentPage('journal'))
                                                                           } } />
        <Route path="data" component={ DataPageContainer } onEnter={ (state, replace) => {
                                                                         store.dispatch(actions.selectPlaybackMode('pause'))
                                                                         store.dispatch(actions.setCurrentPage('data'))
                                                                     } } />
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
