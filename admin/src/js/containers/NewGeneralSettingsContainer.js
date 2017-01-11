
import { connect } from 'react-redux'
import NewGeneralSettingsSection from '../components/NewGeneralSettingsSection'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.getIn(['expeditions', currentExpeditionID])

  return {
    ...ownProps,
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
    }
    // updateExpedition (expedition) {
    //   return dispatch(actions.updateExpedition(expedition))
    // },
  }
}

const NewGeneralSettingsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewGeneralSettingsSection)

export default NewGeneralSettingsContainer
