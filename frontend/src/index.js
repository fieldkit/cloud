// @flow

import React from 'react';
import ReactDOM from 'react-dom';
import App from './js/App';
import './css/index.css';

import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import { createLogger } from 'redux-logger';
import createSagaMiddleware from 'redux-saga';

import { rootSaga } from './js/actions/sagas';
import { isVerboseActionType } from './js/actions/types';
import reducer from './js/reducers';

const sagaMiddleware = createSagaMiddleware();

const loggerMiddleware = createLogger({
    predicate: (getState, action) => !isVerboseActionType(action.type),
    collapsed: (getState, action) => true,
    stateTransformer: state => {
        return state;
    },
});

const createStoreWithMiddleware = applyMiddleware(
    thunkMiddleware,
    sagaMiddleware,
    loggerMiddleware,
)(createStore);

const store = createStoreWithMiddleware(reducer);

sagaMiddleware.run(rootSaga);

ReactDOM.render(<App store={store} />, document.getElementById('root'));
