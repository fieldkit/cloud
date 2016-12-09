
import { connect } from 'react-redux'
import AdminPage from '../components/AdminPage'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {
    expeditions: state.expeditions.get('expeditions'),
    modal: state.expeditions.get('modal')
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    saveChangesAndResume () {
      return dispatch(actions.saveChangesAndResume())
    },
    cancelAction () {
      return dispatch(actions.cancelAction())
    }
  }
}

const AdminPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(AdminPage)

export default AdminPageContainer
