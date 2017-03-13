
import { connect } from 'react-redux'
import UploaderSection from '../../../components/AdminPage/ExpeditionPage/UploaderSection'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {
  const expeditions = state.expeditions
  const projectID = expeditions.getIn(['currentProject', 'id'])
  const expedition = expeditions.get('currentExpedition')

  return {
    ...ownProps,
    projectID,
    expedition,
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setExpeditionProperty (keyPath, value) {
      return dispatch(actions.setExpeditionProperty(keyPath, value))
    }
  }
}

const UploaderContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(UploaderSection)

export default UploaderContainer
