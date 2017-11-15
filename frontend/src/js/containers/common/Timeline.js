
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import Timeline from '../../components/common/Timeline/Timeline'

const mapStateToProps = (state, ownProps) => {
    return {
        ...createSelector(
            state => state.expeditions.get('currentDate'),
            state => state.expeditions.getIn(['expeditions', state.expeditions.get('currentExpedition'), 'startDate']),
            state => state.expeditions.getIn(['expeditions', state.expeditions.get('currentExpedition'), 'endDate']),
            state => state.expeditions.get('documents'),
            (currentDate, startDate, endDate, documents) => ({
                currentDate,
                startDate,
                endDate,
                documents
            })
        )(state)
    }
}

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        updateDate(date, playbackMode, forceUpdate) {
            return dispatch(actions.updateDate(date, playbackMode, forceUpdate))
        }
    }
}

const TimelineContainer = connect(
    mapStateToProps,
    mapDispatchToProps
)(Timeline)

export default TimelineContainer
