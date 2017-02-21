
import { connect } from 'react-redux'
import AdminPage from '../../components/AdminPage/AdminPage'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const projects = state.expeditions.get('projects')
  const expeditions = state.expeditions.get('expeditions')
  const modal = state.expeditions.get('modal')
  const errors = state.expeditions.get('errors')
  const breadcrumbs = state.expeditions.get('breadcrumbs')
  const user = state.expeditions.get('user')

  return {
    expeditions,
    projects,
    modal,
    errors,
    breadcrumbs,
    user
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    saveChangesAndResume () {
      return dispatch(actions.saveChangesAndResume())
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    },
    requestSignOut () {
      return dispatch(actions.requestSignOut())
    }
  }
}

const AdminPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(AdminPage)

export default AdminPageContainer
