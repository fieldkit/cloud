import { combineReducers } from 'redux';

import expeditionReducer from './expeditions'

export default combineReducers({
    expeditions: expeditionReducer
});
