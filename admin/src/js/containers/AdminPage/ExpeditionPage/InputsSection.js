
import { connect } from 'react-redux'
import NewInputsSection from '../../../components/AdminPage/ExpeditionPage/InputsSection'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {
  const expeditions = state.expeditions
  const projectID = expeditions.getIn(['currentProject', 'id'])
  const expedition = expeditions.get('currentExpedition')
  const documentTypes = expedition.get('documentTypes')
    .map(t => {
      return expeditions.getIn(['documentTypes', t.get('id')]) 
    })

  return {
    ...ownProps,
    projectID,
    expedition,
    documentTypes
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setExpeditionProperty (keyPath, value) {
      return dispatch(actions.setExpeditionProperty(keyPath, value))
    },
    fetchSuggestedDocumentTypes (input, type, callback) {
      return dispatch(actions.fetchSuggestedDocumentTypes(input, type, callback))
    },
    addDocumentType (id, type) {
      return dispatch(actions.addDocumentType(id, type))
    },
    removeDocumentType (id) {
      return dispatch(actions.removeDocumentType(id))
    },
    submitInputs () {
      return dispatch(actions.submitInputs())
    }
  }
}

const NewInputsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewInputsSection)

export default NewInputsContainer
