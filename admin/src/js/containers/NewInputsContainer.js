
import { connect } from 'react-redux'
import NewInputsSection from '../components/NewInputsSection'
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
    // },
  }
}

const NewInputsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewInputsSection)

export default NewInputsContainer
