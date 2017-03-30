
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import JournalPage from '../../components/JournalPage/JournalPage'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('expeditions'),
      state => state.expeditions.get('projects'),
      state => state.expeditions.get('documents'),
      state => state.expeditions.get('currentDocuments'),
      state => state.expeditions.get('currentExpedition'),
      state => state.expeditions.get('currentDate'),
      state => state.expeditions.get('forceDateUpdate'),
      (expeditions, projects, documents, currentDocuments, currentExpeditionID, currentDate, forceDateUpdate) => ({
        expeditions,
        projects,
        documents: currentDocuments
          .map(id => documents.get(id))
          .sortBy(d => d.get('date'))
          .reverse(),
        currentExpeditionID,
        currentDate,
        forceDateUpdate
      })
    )(state)
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
