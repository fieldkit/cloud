import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link, browserHistory } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'
import slug from 'slug'
import { protocol, hostname } from '../../../constants/APIBaseURL'

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
      expeditions,
      expedition,
      errors,
      setExpeditionProperty,
      setExpeditionPreset,
      saveGeneralSettings
    } = this.props

    const expeditionName = expedition.get('name')

    return (
      <div id="teams-section" className="section">
        
        <h1 className="section_title">Create a New Expedition <span>(1/2)</span></h1>
        
        { 
          !!expeditions && expeditions.size === 0 &&
          <div className="section_notification">
            <span className="section_notification_icon"></span>
            This project doesn't seem to have any expedition yet. Let's create one.
          </div>
        }

        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <div className="section_form_group">
          <p className="section_form_label">
            Pick a name for this new expedition:
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
                  (!!expeditionName && expeditionName.toLowerCase() !== 'expedition name' ? '' : ' default')
                }
                value={expeditionName}
                onFocus={(e) => {
                  if (!expeditionName || expeditionName === 'Expedition Name') {
                    setExpeditionProperty(['name'], '')
                  }
                }}
                onChange={(e) => {
                  setExpeditionProperty(['name'], e.target.value)
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

        <p className="input-label">
          Your expedition will be available at the following address:
        </p>

        <p className="intro">
          {(protocol + projectID + '.' +  hostname + '/' + slug(expedition.get('name'))).toLowerCase()}
        </p>

        <a href="#" onClick={(e) => {
          e.preventDefault()
          if (!!expeditionName && expeditionName.toLowerCase() !== 'Expedition Name') {
            saveGeneralSettings(() => {
              browserHistory.push('/admin/' + projectID + '/new-expedition/inputs')
            })
          }
        }}>
          <div className={'button hero' + (!expeditionName || expeditionName === 'Expedition Name' ? ' disabled' : '')}>
            Next step: configuring data inputs
          </div>
        </a>

        {
          !!errors &&
          <p className="errors">
            We found one or multiple errors. Please check your information above or try again later.
          </p>
        }

      </div>

    )

  }
}

NewGeneralSettingsSection.propTypes = {
}

export default NewGeneralSettingsSection
