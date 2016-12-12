
import { connect } from 'react-redux'
import NewOutputsSection from '../components/NewOutputsSection'
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

const NewOutputsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewOutputsSection)

export default NewOutputsContainer
