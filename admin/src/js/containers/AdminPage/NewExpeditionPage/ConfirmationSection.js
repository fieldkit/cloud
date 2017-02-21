
import { connect } from 'react-redux'
import NewConfirmationSection from '../../../components/AdminPage/NewExpeditionPage/ConfirmationSection'
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
    saveExpedition () {
      return dispatch(actions.saveExpedition())
    },
  }
}

const NewConfirmationContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewConfirmationSection)

export default NewConfirmationContainer
