
import { connect } from 'react-redux'
import TeamsSection from '../components/TeamsSection'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  const { children, params, disconnect, location } = ownProps

  const expedition = state.expeditions
    .get('expeditions')
    .find(e => {
      return e.get('id') === params.expeditionID
    })

  return {
    ...ownProps,
    expedition
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    updateExpedition (expedition) {
      return dispatch(actions.updateExpedition(expedition))
    }
    // connect: () => {
    // connect () {
    //   return dispatch(actions.connect())
    // },
    // disconnect () {
    //   return dispatch(actions.disconnect())
    // }
  }
}

const TeamsSectionContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(TeamsSection)

export default TeamsSectionContainer
