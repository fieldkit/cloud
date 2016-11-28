
import { connect } from 'react-redux'
import Root from '../components/Root'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {}
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    // connect: () => {
    connect () {
      return dispatch(actions.connect())
    },
    disconnect () {
      return dispatch(actions.disconnect())
    }
  }
}

const RootContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Root)

export default RootContainer
