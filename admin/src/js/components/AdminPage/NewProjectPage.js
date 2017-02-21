import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'
import slug from 'slug'

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
      projects,
      project,
      errors,
      setProjectProperty,
      saveProject
    } = this.props

    const projectName = project.get('name')

    return (
      <div id="teams-section" className="section">
        
        <h1 className="section_title">Create a New Project</h1>

        { 
          !!projects && projects.size === 0 &&
          <div className="section_notification">
            <span className="section_notification_icon"></span>
            Before getting started, we need you to create your first project.
          </div>
        }

        <div className="section_form_group">
          <p className="section_form_label">
            Pick a name for your project:
          </p>
          <div className="section_columns">
            <div
              style={{
                flexBasis: 'calc(66.66% - 1.5625vw)'
              }}
            >
              <input 
                type="text"
                className={
                  'section_form_text_lg' +
                  (!!projectName && projectName.toLowerCase() !== 'project name' ? '' : ' default')
                }
                value={projectName}
                onFocus={(e) => {
                  if (projectName === 'Project Name') {
                    setProjectProperty(['name'], '')
                  }
                }}
                onChange={(e) => {
                  setProjectProperty(['name'], e.target.value)
                }}
              />
              <p className="error"></p>
            </div>
            <div
              className="section_form_description"
              style={{
                flexBasis: 'calc(33.33% - 1.5625vw)'
              }}
            >
              <div>You can change this later if you want.</div>
            </div>
          </div>
          {
            !!errors && !!errors.get && !!errors.get('name') &&
            errors.get('name').map((error, i) => {
              return (
                <p 
                  key={'errors-name-' + i}
                  className="errors"
                >
                  {error}
                </p>
              )
            })
          }
        </div>

        <a href="#" onClick={(e) => {
          e.preventDefault()
          if (!!projectName && projectName.toLowerCase() !== 'Project Name') {
            saveProject(slug(projectName), projectName)
          }
        }}>
          <div className={'button hero ' + (!!projectName && projectName.toLowerCase() !== 'project name' ? '' : 'disabled')}>
            Save this project
          </div>
        </a>

        {
          !!errors &&
          <p className="errors">
            We found one or multiple errors. Please check the information above or try again later.
          </p>
        }

      </div>
    )

  }
}

NewProjectSection.propTypes = {
}

export default NewProjectSection
