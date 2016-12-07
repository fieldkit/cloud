
import { connect } from 'react-redux'
import SignUpPage from '../components/SignUpPage'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {
    errorMessage: state.auth.signUpError,
    success: state.auth.signUpSuccess,
    fetching: state.auth.signUpFetching
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    requestSignUp (email, userName, firstName, lastName, password) {
      dispatch(actions.requestSignUp(email, userName, firstName, lastName, password))
    }
  }
}

const SignUpPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(SignUpPage)

export default SignUpPageContainer
