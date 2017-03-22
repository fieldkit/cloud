
import { connect } from 'react-redux'
import Timeline from '../../components/common/Timeline'
import * as actions from '../../actions'

const mapStateToProps = (state, ownProps) => {

  const currentDate = state.expeditions.get('currentDate')
  const currentExpedition = state.expeditions.get('currentExpedition')
  const startDate = state.expeditions.getIn(['expeditions', currentExpedition, 'startDate'])
  const endDate = state.expeditions.getIn(['expeditions', currentExpedition, 'endDate'])
  const documents = state.expeditions.get('documents')

  return {
    currentDate,
    startDate,
    endDate,
    documents
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    updateDate (date, playbackMode, forceUpdate) {
      return dispatch(actions.updateDate(date, playbackMode, forceUpdate))
    }
  }
}

const TimelineContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Timeline)

export default TimelineContainer
