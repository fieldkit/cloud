
import { connect } from 'react-redux'
import LandingPage from '../components/LandingPage'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    requestSignUp (email, username, password, invite, project) {
      dispatch(actions.requestSignUp(email, username, password, invite, password))
    }
  }
}

const LandingPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(LandingPage)

export default LandingPageContainer
