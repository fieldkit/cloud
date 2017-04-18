
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import NotificationPanel from '../../components/MapPage/NotificationPanel'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('currentDate'),
      state => state.expeditions.get('documents'),
      state => state.expeditions.get('currentDocuments'),
      state => state.expeditions.get('showSensors'),
      (currentDate, documents, currentDocuments, showSensors) => ({
        currentDocuments: documents
          .filter(d => state.expeditions.get('currentDocuments').includes(d.get('id')))
          .filter(d => Math.abs(d.get('date') - currentDate + 100000) < 200000)
          .filter(d => showSensors ? true : ! d.has('GPSSpeed') )
          .sortBy(d => d.get('date'))
          .slice(0,5)
      })
    )(state)
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {}
}

const NotificationPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NotificationPanel)

export default NotificationPanelContainer
