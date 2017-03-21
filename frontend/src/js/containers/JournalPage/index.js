
import { connect } from 'react-redux'
import JournalPage from '../../components/JournalPage'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const currentExpeditionID = state.expeditions.get('currentExpedition')
  const projects = state.expeditions.get('projects')
  const expeditions = state.expeditions.get('expeditions')
  const documents = state.expeditions.get('currentDocuments')
    .map(id => state.expeditions.getIn(['documents', id]))

  return {
    expeditions,
    projects,
    documents,
    currentExpeditionID
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    updateDate (date, playbackMode) {
      return dispatch(actions.updateDate(date, playbackMode))
    }
  }
}

const JournalPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(JournalPage)

export default JournalPageContainer
