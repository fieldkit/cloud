
import { connect } from 'react-redux'
import JournalPage from '../../components/JournalPage'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const projects = state.expeditions.get('projects')
  const expeditions = state.expeditions.get('expeditions')
  const modal = state.expeditions.get('modal')

  return {
    expeditions,
    projects,
    modal
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

const JournalPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(JournalPage)

export default JournalPageContainer
