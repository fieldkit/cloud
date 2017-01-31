import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link, browserHistory } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'
import slug from 'slug'

class NewGeneralSettingsSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      addMemberValue: null,
      inputValues: {}
    }
  }

  render () {
    const {
      projectID,
      expedition,
      setExpeditionProperty,
      setExpeditionPreset,
      saveGeneralSettings,
      error
    } = this.props

    return (
      <div id="teams-section" className="section">
        <div className="section-header">
          <h1>Create a New Expedition (1/2)</h1>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <p className="input-label">
          Pick a name for your expedition:
        </p>
        <div className="columns-container">
          <div className="main-input-container">
            <input 
              type="text"
              value={expedition.get('name')}
              onFocus={(e) => {
                if (expedition.get('name') === 'New Expedition') {
                  setExpeditionProperty(['name'], '')
                }
              }}
              onChange={(e) => {
                setExpeditionProperty(['name'], e.target.value)
              }}
            />
            {
              !!error && !!error.get('name') &&
              error.get('name').map((error, i) => {
                return (
                  <p 
                    key={'error-' + i}
                    className="error"
                  >
                    {error}
                  </p>
                )
              })
            }
          </div>
          <p className="input-description">
            Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
          </p>
        </div>

        <p className="input-label">
          Your expedition will be available at the following address:
        </p>

        <p className="intro">
          {'https://' + projectID + '.fieldkit.org/' + slug(expedition.get('name'))}
        </p>

        {
          !!error &&
          <p className="error">
            We found one or multiple errors. Please check your information above or try again later.
          </p>
        }

        <a href="#" onClick={(e) => {
          e.preventDefault()
          saveGeneralSettings(() => {
            browserHistory.push('/admin/' + projectID + '/new-expedition/inputs')
          })
        }}>
          <div className="button hero">
            Next step: configuring data inputs
          </div>
        </a>

      </div>
    )

  }
}

NewGeneralSettingsSection.propTypes = {
}

export default NewGeneralSettingsSection
