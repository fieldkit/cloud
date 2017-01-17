
import { connect } from 'react-redux'
import NewProjectSection from '../components/NewProjectSection'
import * as actions from '../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const currentProjectID = expeditions.get('currentProjectID')
  const currentProject = expeditions.getIn(['projects', currentProjectID])

  return {
    ...ownProps,
    currentProjectID,
    currentProject,
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setProjectProperty (key, value) {
      return dispatch(actions.setProjectProperty(key, value))
    },
    createProject (name) {
      return dispatch(actions.createProject(name))
    }
  }
}

const NewProjectContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewProjectSection)

export default NewProjectContainer
