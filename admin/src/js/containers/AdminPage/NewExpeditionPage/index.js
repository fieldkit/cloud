
import { connect } from 'react-redux'
import NewExpeditionPage from '../../../components/AdminPage/NewExpeditionPage'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentProjectID = expeditions.get('currentProjectID')
  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.getIn(['expeditions', currentExpeditionID])

  return {
    ...ownProps,
    currentProjectID,
    currentExpedition
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setExpeditionProperty (key, value) {
      return dispatch(actions.setExpeditionProperty(key, value))
    },
    setExpeditionPreset (type) {
      return dispatch(actions.setExpeditionPreset(type))
    },
    submitGeneralSettings () {
      return dispatch(actions.submitGeneralSettings())
    }
    // updateExpedition (expedition) {
    //   return dispatch(actions.updateExpedition(expedition))
    // },
  }
}

const NewGeneralSettingsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewExpeditionPage)

export default NewGeneralSettingsContainer
