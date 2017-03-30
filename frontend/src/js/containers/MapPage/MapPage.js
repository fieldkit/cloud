
import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import MapPage from '../../components/MapPage/MapPage'


const mapStateToProps = (state, ownProps) => {
  return {
    ...createSelector(
      state => state.expeditions.get('documents'),
      (documents) => ({
        documents
      })
    )(state)
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

const MapPageContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(MapPage)

export default MapPageContainer
