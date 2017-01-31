
import { connect } from 'react-redux'
import GeneralSettingsSection from '../../../components/AdminPage/ExpeditionPage/GeneralSettingsSection'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const projectID = expeditions.get('currentProjectID')
  const expedition = expeditions.get('newExpedition')

  return {
    ...ownProps,
    projectID,
    expedition
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
    saveGeneralSettings () {
      return dispatch(actions.saveGeneralSettings())
    }
  }
}

const GeneralSettingsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(GeneralSettingsSection)

export default GeneralSettingsContainer
