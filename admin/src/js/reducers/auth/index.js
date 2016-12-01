
import * as actions from '../../actions'
import I from 'immutable'

export const initialState = I.fromJS({
  signInError: null,
  signInFetching: false,
  signUpError: null,
  signUpSuccess: false,
  signUpFetching: false
})

const authReducer = (state = initialState, action) => {
  switch (action.type) {
    case actions.LOGIN_REQUEST:
      return state
        .set('signInFetching', true)
        .set('signInError', null)
    case actions.LOGIN_SUCCESS:
      return state
        .set('signInFetching', false)
        .set('signInError', null)
    case actions.LOGIN_ERROR:
      return state
        .set('signInFetching', false)
        .set('signInError', action.message)
    case actions.SIGNUP_REQUEST:
      return state
        .set('signUpFetching', true)
        .set('signUpSuccess', false)
        .set('signUpError', null)
    case actions.SIGNUP_SUCCESS:
      return state
        .set('signUpFetching', false)
        .set('signUpSuccess', true)
        .set('signUpError', null)
    case actions.SIGNUP_ERROR:
      return state
        .set('signUpFetching', false)
        .set('signUpSuccess', false)
        .set('signUpError', action.message)
    default: 
      return state
  }
}


export default authReducer
