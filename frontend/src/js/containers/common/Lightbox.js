
import { connect } from 'react-redux'
import * as actions from '../../actions'
import Lightbox from '../../components/common/Lightbox'

const mapStateToProps = (state, ownProps) => {
  const lightboxDocumentID = state.expeditions.get('lightboxDocumentID')
  const data = state.expeditions.getIn(['documents', lightboxDocumentID])
  const currentExpeditionID = state.get('currentExpedition')
  return {
    data,
    currentExpeditionID
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