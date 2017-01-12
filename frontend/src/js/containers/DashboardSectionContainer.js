
import { connect } from 'react-redux'
import DashboardSection from '../components/DashboardSection'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.getIn(['expeditions', currentExpeditionID])

  console.log(state.expeditions.toJS(), currentExpeditionID, currentExpedition)

  return {
    ...ownProps,
    currentExpedition
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setExpeditionProperty (key, value) {
      return dispatch(actions.setExpeditionProperty(key, value))
    }
    // updateExpedition (expedition) {
    //   return dispatch(actions.updateExpedition(expedition))
    // },
  }
}

const DashboardSectionContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(DashboardSection)

export default DashboardSectionContainer
