
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import Lightbox from '../../components/common/Lightbox'

const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.getIn(['documents', state.expeditions.get('lightboxDocumentID') || 'null']),
      state => state.expeditions.get('currentExpedition'),
      (data, currentExpeditionID) => ({
        data,
        currentExpeditionID
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
    }
  }
}

const LightboxContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Lightbox)

export default LightboxContainer