
import { connect } from 'react-redux'
import Root from '../components/Root'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  return {
    foo: 'bar'
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    // initialize: (currentDate) => {
    //   return dispatch(actions.initialize(null))
    // }
  }
}

const RootContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Root)

export default RootContainer
