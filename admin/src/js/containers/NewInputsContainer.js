
import { connect } from 'react-redux'
import NewInputsSection from '../components/NewInputsSection'
import * as actions from '../actions'


const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentProjectID = expeditions.get('currentProjectID')
  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.getIn(['expeditions', currentExpeditionID])
  const currentDocumentTypeID = expeditions.get('currentDocumentTypeID')
  const documentTypes = expeditions.getIn(['expeditions', currentExpeditionID, 'documentTypes'])
    .map(t => {
      return expeditions.getIn(['documentTypes', t.get('id')]) 
    })

  // const currentDocumentType = !!documentTypes.size ? documentTypes.find(t => {
  //   // console.log('wow', t)
  //   return t.get('id') === currentDocumentTypeID
  // }) : null

  return {
    ...ownProps,
    currentProjectID,
    currentExpedition,
    documentTypes,
    // currentDocumentType
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
    }
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
