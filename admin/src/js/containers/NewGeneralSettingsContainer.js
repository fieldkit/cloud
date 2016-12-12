
import { connect } from 'react-redux'
import NewGeneralSettingsSection from '../components/NewGeneralSettingsSection'
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

const NewGeneralSettingsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewGeneralSettingsSection)

export default NewGeneralSettingsContainer
