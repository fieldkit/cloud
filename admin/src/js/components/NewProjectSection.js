import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
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
      currentExpedition,
      setExpeditionProperty,
      setExpeditionPreset
    } = this.props

    return (
      <div id="teams-section" className="section">
        <div className="section-header">
          <h1>Create a New Expedition</h1>
          <p>process breadcrumbs</p>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <p className="input-label">
          Pick a name for your expedition:
        </p>
        <div className="columns-container">
          <div className="main-input-container">
            <input type="text" value={currentExpedition.get('name')} onChange={(e) => {
              setExpeditionProperty(['name'], e.target.value)
            }}/>
            <p className="error"></p>
          </div>
          <p className="input-description">
            Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
          </p>
        </div>

        <p className="input-label">
          Pick a URL for the visualization page:
        </p>
        <div className="columns-container">
          <div className="main-input-container">
            <input type="text" value={'fieldkit.org/' + currentExpedition.get('id')} onChange={(e) => {
              setExpeditionProperty(['id'], e.target.value.slice(13))
            }}/>
            <p className="error"></p>
          </div>
          <p className="input-description">
            Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
          </p>
        </div>

        <p className="input-label">
          Describe your project:
        </p>
        <div className="columns-container">
          <div className="main-input-container">
            <textarea
              rows="4"
              cols="50"
              onChange={e => {
                setExpeditionProperty(['description'], e.target.value)
              }}
              value={currentExpedition.get('description')}
            >
            </textarea>
          </div>
          <p className="input-description">
            Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
          </p>
        </div>

        <p className="input-label">
          Do you have an idea in mind? You can select one of these presets:
        </p>
        <div className="columns-container">
          <ul className="main-input-container">
            <li
              className={ currentExpedition.get('selectedPreset') === 'rookie' ? 'selected' : ''}
              onClick={() => {
                setExpeditionPreset('rookie')
              }}
            >
              <h2>
                Rookie expedition
              </h2>
              <p>
                Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
              </p>
            </li>
            <li
              className={ currentExpedition.get('selectedPreset') === 'advanced' ? 'selected' : ''}
              onClick={() => {
                setExpeditionPreset('advanced')
              }}
            >
              <h2>
                Advanced expedition
              </h2>
              <p>
                Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
              </p>
            </li>
            <li
              className={ currentExpedition.get('selectedPreset') === 'pro' ? 'selected' : ''}
              onClick={() => {
                setExpeditionPreset('pro')
              }}
            >
              <h2>
                Start from scratch
              </h2>
              <p>
                Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
              </p>
            </li>
          </ul>
          <p className="input-description">
            Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien.
          </p>
        </div>

        <p className="status">
        </p>

        <Link to={'/admin/' + currentProjectID + '/new-expedition/inputs'}>
          <div className="button hero">
            Next step: configuring document types
          </div>
        </Link>

      </div>
    )

  }
}

NewProjectSection.propTypes = {
}

export default NewProjectSection
