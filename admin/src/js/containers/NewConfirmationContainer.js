
import { connect } from 'react-redux'
import NewConfirmationSection from '../components/NewConfirmationSection'
import * as actions from '../actions'


const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentProjectID = expeditions.get('currentProjectID')
  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.getIn(['expeditions', currentExpeditionID])

  return {
    ...ownProps,
    currentProjectID,
    currentExpedition,
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    // updateExpedition (expedition) {
    //   return dispatch(actions.updateExpedition(expedition))
    // },
  }
}

const NewConfirmationContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewConfirmationSection)

export default NewConfirmationContainer
