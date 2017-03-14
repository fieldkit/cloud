
import { connect } from 'react-redux'
import TeamsSection from '../../../components/AdminPage/ExpeditionPage/TeamsSection'
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

const NewInputsContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(TeamsSection)

export default NewInputsContainer
