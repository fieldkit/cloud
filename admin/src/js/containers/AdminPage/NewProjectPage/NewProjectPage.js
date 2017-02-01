
import { connect } from 'react-redux'
import NewProjectPage from '../../../components/AdminPage/NewProjectPage'
import * as actions from '../../../actions'

const mapStateToProps = (state, ownProps) => {

  const expeditions = state.expeditions
  const projects = expeditions.get('projects')
  const project = expeditions.get('newProject')

  return {
    ...ownProps,
    projects,
    project
  }
}

const mapDispatchToProps = (dispatch, ownProps, state) => {
  return {
    setProjectProperty (key, value) {
      return dispatch(actions.setProjectProperty(key, value))
    },
    saveProject (id, name) {
      return dispatch(actions.saveProject(id, name))
    }
  }
}

const NewProjectContainer = connect(
  mapStateToProps,
  mapDispatchToProps
)(NewProjectPage)

export default NewProjectContainer
