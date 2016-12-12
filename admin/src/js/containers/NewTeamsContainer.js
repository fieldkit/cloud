
import { connect } from 'react-redux'
import NewTeamsSection from '../components/NewTeamsSection'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  
  return {
    ...ownProps
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    // updateExpedition (expedition) {
    //   return dispatch(actions.updateExpedition(expedition))
    // }
  }
}

const NewTeamsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewTeamsSection)

export default NewTeamsContainer
