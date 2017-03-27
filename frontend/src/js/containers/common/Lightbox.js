
import { connect } from 'react-redux'
import * as actions from '../../actions'
import Lightbox from '../../components/common/Lightbox'

const mapStateToProps = (state, ownProps) => {
  const lightboxDocumentID = state.expeditions.get('lightboxDocumentID')
  const data = !!lightboxDocumentID ? state.expeditions.getIn(['documents', lightboxDocumentID]) : null
  const currentExpeditionID = state.expeditions.get('currentExpedition')
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