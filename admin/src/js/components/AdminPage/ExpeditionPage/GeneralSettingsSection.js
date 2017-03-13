import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'
import slug from 'slug'
import { protocol, hostname } from '../../../constants/APIBaseURL'

class GeneralSettingsSection extends React.Component {
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
      <div className="section">
        <h1 className="section_title">General Settings</h1>

        <p className="intro">
          Configure this expedition's main settings. Changes are automatically saved.
        </p>

        <div className="section_form_group">
          <p className="section_form_label">
            Edit this expedition's name:
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

        <div className="status">
          Changes saved!
        </div>

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

GeneralSettingsSection.propTypes = {
}

export default GeneralSettingsSection
