
import { connect } from 'react-redux'
import InputsSection from '../../../components/AdminPage/ExpeditionPage/InputsSection'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {
  const expeditions = state.expeditions
  const projectID = expeditions.getIn(['currentProject', 'id'])
  const expedition = expeditions.get('currentExpedition')
  const inputs = expeditions.get('inputs')
    .filter(t => {
      return expedition.get('inputs').includes(t.get('id'))
    })

  return {
    ...ownProps,
    projectID,
    expedition,
    inputs
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setExpeditionProperty (keyPath, value) {
      return dispatch(actions.setExpeditionProperty(keyPath, value))
    },
    fetchSuggestedInputs (input, type, callback) {
      return dispatch(actions.fetchSuggestedInputs(input, type, callback))
    },
    addInput (id, type) {
      return dispatch(actions.addInput(id, type))
    },
    removeInput (id) {
      return dispatch(actions.removeInput(id))
    },
    submitInputs () {
      return dispatch(actions.submitInputs())
    }
  }
}

const InputsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(InputsSection)

export default InputsContainer
