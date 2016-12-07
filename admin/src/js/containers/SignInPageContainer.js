
import { connect } from 'react-redux'
import SignInPage from '../components/SignInPage'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {
    errorMessage: state.auth.signInError,
    fetching: state.auth.signInFetching,
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    connect () {
      return dispatch(actions.connect())
    },
    requestSignIn (userName, password) {
      dispatch(actions.requestSignIn(userName, password))
    }
  }
}

const SignInPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(SignInPage)

export default SignInPageContainer
