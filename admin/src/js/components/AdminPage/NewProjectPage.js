import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'

class NewProjectSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      addMemberValue: null,
      inputValues: {}
    }
  }

  render () {

    const { 
      currentProjectID,
      currentProject,
      setProjectProperty,
      createProject
    } = this.props

    return (
      <div id="teams-section" className="section">
        <div className="section-header">
          <h1>Create a New Project</h1>
        </div>

        <p className="input-label">
          Pick a name for your project:
        </p>
        <div className="columns-container">
          <div className="main-input-container">
            <input 
              type="text"
              value={!!currentProject ? currentProject.get('name') : 'New Project'}
              onFocus={(e) => {
                if (!!currentProject ? currentProject.get('name') === 'New Project' : 'New Project') {
                  setProjectProperty(['name'], '')
                }
              }}
              onChange={(e) => {
                setProjectProperty(['name'], e.target.value)
              }}
            />
            <p className="error"></p>
          </div>
          <p className="input-description">
            You will be able to change this name later if necessary.
          </p>
        </div>

        <p className="status">
        </p>

        <a href="#" onClick={(e) => {
          e.preventDefault()
          createProject(currentProject.get('name'))
        }}>
          <div className="button hero">
            Create this project
          </div>
        </a>

      </div>
    )

  }
}

NewProjectSection.propTypes = {
}

export default NewProjectSection
