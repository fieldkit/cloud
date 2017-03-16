
import { connect } from 'react-redux'
import LandingPage from '../../components/LandingPage/LandingPage'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {
  const errors = state.expeditions.get('errors')
  return {
    errors
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    requestSignUp (email, username, password, invite) {
      dispatch(actions.requestSignUp(email, username, password, invite))
    },
    requestSignIn (username, password) {
      dispatch(actions.requestSignIn(username, password))
    }
  }
}

const LandingPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(LandingPage)

export default LandingPageContainer
