import { connect } from 'react-redux'
import * as actions from '../../actions'
import { createSelector } from 'reselect'

import DisclaimerPanel from '../../components/MapPage/DisclaimerPanel'

const mapStateToProps = (state, ownProps) => {
  return {}
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {}
}

const DisclaimerPanelContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DisclaimerPanel)

export default DisclaimerPanelContainer
