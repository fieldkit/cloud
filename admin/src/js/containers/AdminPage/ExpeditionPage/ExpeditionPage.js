
import { connect } from 'react-redux'
import DashboardSection from '../../../components/AdminPage/ExpeditionPage/ExpeditionPage'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  // const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentExpedition = expeditions.get('newExpedition')

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
