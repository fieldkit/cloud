
import { connect } from 'react-redux'
import NewOutputsSection from '../../../components/NewOutputsSection'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.getIn(['expeditions', currentExpeditionID])

  console.log(currentExpeditionID)

  return {
    ...ownProps,
    currentExpedition
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setExpeditionProperty (keyPath, value) {
      return dispatch(actions.setExpeditionProperty(keyPath, value))
    },
    // updateExpedition (expedition) {
    //   return dispatch(actions.updateExpedition(expedition))
    // },
  }
}

const NewOutputsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewOutputsSection)

export default NewOutputsContainer
