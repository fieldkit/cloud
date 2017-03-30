
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import Lightbox from '../../components/common/Lightbox'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.getIn(['documents', state.expeditions.get('lightboxDocumentID') || 'null']),
      state => state.expeditions.get('currentExpedition'),
      state => state.expeditions.get('previousDocumentID'),
      state => state.expeditions.get('nextDocumentID'),
      (data, currentExpeditionID, previousDocumentID, nextDocumentID) => ({
        data,
        currentExpeditionID, 
        previousDocumentID,
        nextDocumentID
      })
    )(state)
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    updateDate (date, playbackMode) {
      return dispatch(actions.updateDate(date, playbackMode))
    },
    closeLightbox () {
      return dispatch(actions.closeLightbox())
    },
    openLightbox (id) {
      return dispatch(actions.openLightbox(id))
    }
  }
}

const LightboxContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Lightbox)

export default LightboxContainer