
import { connect } from 'react-redux'
import TeamsSection from '../components/TeamsSection'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {
  const expeditions = state.expeditions

  const currentExpeditionID = expeditions.get('currentExpeditionID')
  const currentTeamID = expeditions.get('currentTeamID')

  const expedition = expeditions.getIn(['expeditions', currentExpeditionID])

  const teams = expeditions.getIn(['expeditions', currentExpeditionID, 'teams'])
    .map(t => {
      return expeditions.getIn(['teams', t]) 
    })

  const members = teams.size > 0 ? (expeditions.getIn(['teams', currentTeamID, 'members'])
    .map(m => {
      return expeditions.getIn(['people', m.get('id')])
    })) : []

  const currentTeam = !!teams.size ? teams.find(t => {
    return t.get('id') === currentTeamID
  }) : null

  const editedTeam = expeditions.get('editedTeam')

  return {
    ...ownProps,
    expedition,
    teams,
    members,
    currentTeam,
    editedTeam,
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    addTeam () {
      return dispatch(actions.addTeam())
    },
    setCurrentTeam (teamID) {
      return dispatch(actions.setCurrentTeam(teamID))
    },
    setCurrentMember (memberID) {
      return dispatch(actions.setCurrentMember(memberID))
    },
    removeCurrentTeam () {
      return dispatch(actions.removeCurrentTeam())
    },
    startEditingTeam () {
      return dispatch(actions.startEditingTeam())
    },
    stopEditingTeam () {
      return dispatch(actions.stopEditingTeam())
    },
    setTeamProperty (key, value) {
      return dispatch(actions.setTeamProperty(key, value))
    },
    setMemberProperty (memberID, key, value) {
      return dispatch(actions.setMemberProperty(memberID, key, value))
    },
    saveChangesToTeam () {
      return dispatch(actions.saveChangesToTeam())
    },
    clearChangesToTeam () {
      return dispatch(actions.clearChangesToTeam())
    },
    fetchSuggestedMembers (input, callback) {
      return dispatch(actions.fetchSuggestedMembers(input, callback))
    },
    addMember (id) {
      return dispatch(actions.addMember(id))
    },
    removeMember (id) {
      return dispatch(actions.removeMember(id))
    }
  }
}

const TeamsSectionContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(TeamsSection)

export default TeamsSectionContainer
